package admin

import (
	"fmt"
	"github.com/zeebo/sign"
	"launchpad.net/gobson/bson"
	"math"
	"net/http"
	"path"
	"reflect"
	"strings"
)

//Parse request grabs the paramaters out of the request URL for the collection
//and object id the handler will operate on.
func parseRequest(p string) (coll, id string) {
	chunks := strings.Split(path.Clean(p), "/")
	if chunks[0] == "." {
		return
	}
	coll = chunks[0]
	if len(chunks) >= 2 {
		id = chunks[1]
	}
	return
}

func (a *Admin) baseContext(req *http.Request) (ctx BaseContext) {
	ctx.Managed = a.index_cache
	ctx.Reverser = Reverser{a}

	if auth, ex := a.auth_cache[req]; ex {
		ctx.Auth = &auth
	}

	return
}

func (a *Admin) auth(w http.ResponseWriter, req *http.Request) {
	action, n := parseRequest(req.URL.Path)

	//ensure that we have an action and nothing else
	if action == "" || n != "" {
		a.Renderer.NotFound(w, req)
		return
	}

	switch action {
	case "logout":
		as := AuthSession{}
		as.clear(w)

		a.Renderer.LoggedOut(w, req, a.baseContext(req))
		return
	case "login":
		//pass down through the switch
	default:
		a.Renderer.NotFound(w, req)
		return
	}

	//do the easy case first
	if req.Method != "POST" {
		a.Renderer.Authorize(w, req, AuthorizeContext{
			BaseContext: a.baseContext(req),
		})
		return
	}

	resp := a.Auth.Authorize(req)

	var success bool

	//gotta set the cookie
	if resp.Passed {
		as := AuthSession{
			Username: resp.Username,
			Key:      resp.Key,
		}
		signer := sign.Signer{a.Key}

		if err := as.add(signer, w); err != nil {
			a.logger.Printf("Error adding auth cookie: %s", err)
			resp.Error = err.Error()
			goto render
		}

		success = true
		//look up the redirect url
		if val, err := req.Cookie("redirect"); err == nil {
			http.Redirect(w, req, val.Value, http.StatusTemporaryRedirect)
			return
		}
	}

render:

	a.Renderer.Authorize(w, req, AuthorizeContext{
		BaseContext: a.baseContext(req),
		Attempted:   true,
		Success:     success,
		Error:       resp.Error,
	})
}

//Presents the detail view for an object in a collection
func (a *Admin) detail(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)

	//ensure we have both a collection key and an id
	if coll == "" || id == "" {
		a.Renderer.NotFound(w, req)
		return
	}

	//make sure we know about the requested collection
	if !a.hasType(coll) {
		a.Renderer.NotFound(w, req)
		return
	}

	c, t := a.collFor(coll), a.newType(coll)

	//load into T
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(t); err != nil {
		if err.Error() == "Document not found" {
			a.Renderer.NotFound(w, req)
			return
		}
		a.Renderer.InternalError(w, req, err)
		return
	}

	//create the values for the template
	ctx, err := generateContext(t, nil)
	if err != nil {
		a.Renderer.InternalError(w, req, err)
		return
	}

	a.Renderer.Detail(w, req, DetailContext{
		BaseContext: a.baseContext(req),
		Collection:  coll,
		Object:      t,
		Form: Form{
			template: a.types[coll].Template,
			context:  ctx,
			logger:   a.logger,
		},
	})
}

//Presents the delete view for an object in a collection
func (a *Admin) delete(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)

	//ensure we have both a collection and an id
	if coll == "" || id == "" {
		a.Renderer.NotFound(w, req)
		return
	}

	//make sure we know about the requested collection
	if !a.hasType(coll) {
		a.Renderer.NotFound(w, req)
		return
	}

	c, t := a.collFor(coll), a.newType(coll)

	//load into T
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(t); err != nil {
		if err.Error() == "Document not found" {
			a.Renderer.NotFound(w, req)
			return
		}
		a.Renderer.InternalError(w, req, err)
		return
	}

	//check if they're sure they want to delete
	req.ParseForm()

	var attempted, success bool
	var err error

	if req.Form.Get("_sure") == "yes" {
		attempted = true

		err = c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
		success = err == nil
	}

	//create the values for the template
	ctx, err := generateContext(t, nil)
	if err != nil {
		a.Renderer.InternalError(w, req, err)
		return
	}

	a.Renderer.Delete(w, req, DeleteContext{
		BaseContext: a.baseContext(req),
		Collection:  coll,
		Object:      t,
		Attempted:   attempted,
		Success:     success,
		Error:       err,
		Form: Form{
			template: a.types[coll].Template,
			context:  ctx,
			logger:   a.logger,
		},
	})
}

//Presents the index page giving an overall view of the database
func (a *Admin) index(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)

	//ensure we have neither a collection nor an id
	if coll != "" || id != "" {
		a.Renderer.NotFound(w, req)
		return
	}

	a.Renderer.Index(w, req, a.baseContext(req))
}

//Presents a list of objects in a collection matching filtering/sorting criteria
func (a *Admin) list(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)

	//ensure we have a collection, but no id
	if coll == "" || id != "" {
		a.Renderer.NotFound(w, req)
	}

	//make sure we know about the requested collection
	if !a.hasType(coll) {
		a.Renderer.NotFound(w, req)
		return
	}

	c := a.collFor(coll)

	//TODO: make this load into a map[string]interface{} instead
	//to reduce the amount of reflection we need to do. We can't get
	//objects that way though so see if thats an issue.

	total, err := c.Find(nil).Count()
	if err != nil {
		a.Renderer.InternalError(w, req, err)
	}

	//grab the data
	var items []interface{}
	iter, page, numpage := listParse(c, req.URL.Query())
	for {
		t := a.newType(coll)
		if !iter.Next(t) {
			break
		}
		items = append(items, t)
	}

	//report any errors our iterator made
	if err := iter.Err(); err != nil {
		a.Renderer.InternalError(w, req, err)
		return
	}

	ids := a.types[coll].ColumnIds

	//generate the columns
	columns := make([]string, len(ids))
	typ := a.types[coll].Type
	for j, idx := range ids {
		columns[j] = typ.Field(idx).Name
	}

	//make the values :(
	values := make([][]string, len(items))
	for i, obj := range items {
		val, err := indirect(reflect.ValueOf(obj))
		if err != nil {
			a.Renderer.InternalError(w, req, err)
			return
		}

		values[i] = make([]string, len(ids))

		for j, idx := range ids {
			var data string

			//TODO: make things that involve hexable into a function that gets
			//called.
			switch item := val.Field(idx).Interface().(type) {
			case hexable:
				data = item.Hex()
			default:
				data = fmt.Sprint(item)
			}

			values[i][j] = data
		}
	}

	pages := int(math.Ceil(float64(total) / float64(numpage)))
	if pages < 1 {
		pages = 1
	}

	a.Renderer.List(w, req, ListContext{
		BaseContext: a.baseContext(req),
		Collection:  coll,
		Columns:     columns,
		Values:      values,
		Objects:     items,
		Pagination: Pagination{
			Pages:       pages,
			CurrentPage: page,
			query:       req.URL.Query(),
		},
	})
}

//Presents a handler that updates an object and shows the results of the update
func (a *Admin) update(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)

	//ensure we have both an id and a collection
	if coll == "" || id == "" {
		a.Renderer.NotFound(w, req)
		return
	}

	//make sure we know about the requested collection
	if !a.hasType(coll) {
		a.Renderer.NotFound(w, req)
		return
	}

	c, t := a.collFor(coll), a.newType(coll)

	//grab the data
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(t); err != nil {
		if err.Error() == "Document not found" {
			a.Renderer.NotFound(w, req)
			return
		}
		a.Renderer.InternalError(w, req, err)
		return
	}

	var attempted, success bool
	var errors map[string]error
	if req.Method == "POST" {
		attempted = true

		var err error
		errors, err = performLoading(req, t)
		if err != nil {
			a.Renderer.InternalError(w, req, err)
			return
		}
		if errors != nil && len(errors) > 0 {
			goto render
		}

		if err := c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, t); err != nil {
			a.Renderer.InternalError(w, req, err)
			return
		}
		success = true
	}

render:
	var form = Form{
		template: a.types[coll].Template,
		logger:   a.logger,
	}
	if ctx, err := generateContext(t, errors); err != nil {
		a.Renderer.InternalError(w, req, err)
		return
	} else {
		form.context = ctx
	}

	a.Renderer.Update(w, req, UpdateContext{
		BaseContext: a.baseContext(req),
		Collection:  coll,
		Object:      t,
		Attempted:   attempted,
		Success:     success,
		Form:        form,
	})
}

//Presents a handler that creates an object and shows the results of the create
func (a *Admin) create(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)

	//ensure we have a collection but no id
	if coll == "" || id != "" {
		a.Renderer.NotFound(w, req)
	}

	//make sure we know about the requested collection
	if !a.hasType(coll) {
		a.Renderer.NotFound(w, req)
		return
	}

	c, t := a.collFor(coll), a.newType(coll)

	var attempted, success bool
	var errors map[string]error
	if req.Method == "POST" {
		attempted = true

		var err error
		errors, err = performLoading(req, t)
		if err != nil {
			a.Renderer.InternalError(w, req, err)
			return
		}
		if errors != nil && len(errors) > 0 {
			goto render
		}

		id, err := c.Upsert(d{"_id": ""}, t)
		if err != nil {
			a.Renderer.InternalError(w, req, err)
			return
		}

		//lets grab the thing back out from the database
		if err = c.Find(bson.M{"_id": id}).One(t); err != nil {
			a.Renderer.InternalError(w, req, err)
			return
		}

		success = true
	}

render:
	var form = Form{
		template: a.types[coll].Template,
		logger:   a.logger,
	}
	if attempted {
		if ctx, err := generateContext(t, errors); err != nil {
			a.Renderer.InternalError(w, req, err)
			return
		} else {
			form.context = ctx
		}
	} else {
		val, err := CreateEmptyValues(t)
		if err != nil {
			a.Renderer.InternalError(w, req, err)
		}

		form.context = TemplateContext{
			Values: val,
			Errors: errors,
		}
	}

	a.Renderer.Create(w, req, CreateContext{
		BaseContext: a.baseContext(req),
		Collection:  coll,
		Attempted:   attempted,
		Success:     success,
		Form:        form,
	})
}

//performLoading is a helper function that loads and validates the form, returning
//any errors from the two steps. It respects if the type is a Loader.
func performLoading(req *http.Request, t Formable) (errors map[string]error, err error) {
	//TODO: files!
	err = req.ParseForm()
	if err != nil {
		return
	}

	if l, ok := t.(Loader); ok {
		errors, err = l.Load(req.Form)
	} else {
		errors, err = Load(req.Form, t)
	}

	//do we have loading errors?
	if (errors != nil && len(errors) > 0) || err != nil {
		return
	}

	errors = t.Validate()
	return
}

//generateContext takes a value that should be filled in, and some errors generated
//while filling it in and returns a TemplateContext for rendering a Form, and
//any errors attempting to do so.
func generateContext(t Formable, errors map[string]error) (TemplateContext, error) {
	if l, ok := t.(Loader); ok {
		return TemplateContext{
			Values: l.GenerateValues(),
			Errors: errors,
		}, nil
	}

	values, err := CreateValues(t)
	if err != nil {
		return TemplateContext{}, err
	}

	return TemplateContext{
		Values: values,
		Errors: errors,
	}, nil
}

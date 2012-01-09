package admin

import (
	"launchpad.net/gobson/bson"
	"net/http"
	"path"
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

//Presents the detail view for an object in a collection
func (a *Admin) Detail(w http.ResponseWriter, req *http.Request) {
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

	a.Renderer.Detail(w, req, DetailContext{
		Object: t,
	})
}

//Presents the index page giving an overall view of the database
func (a *Admin) Index(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)

	//ensure we have neither a collection nor an id
	if coll != "" || id != "" {
		a.Renderer.NotFound(w, req)
		return
	}

	a.generateIndexCache()
	a.Renderer.Index(w, req, IndexContext{
		Managed: a.index_cache,
	})
}

//Presents a list of objects in a collection matching filtering/sorting criteria
func (a *Admin) List(w http.ResponseWriter, req *http.Request) {
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

	//grab the data
	var items []interface{}
	iter := listParse(c, req.URL.Query())
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
	}

	a.Renderer.List(w, req, ListContext{
		Objects: items,
	})
}

//Presents a handler that updates an object and shows the results of the update
func (a *Admin) Update(w http.ResponseWriter, req *http.Request) {
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

	//attempt to update

	a.Renderer.Update(w, req, UpdateContext{
		Object: t,
	})
}

//Presents a handler that creates an object and shows the results of the create
func (a *Admin) Create(w http.ResponseWriter, req *http.Request) {
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

	//attempt to insert

	a.Renderer.Create(w, req, CreateContext{})
}

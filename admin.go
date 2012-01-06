package admin

import (
	"fmt"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"net/http"
	"path"
	"reflect"
	"strings"
)

//useful type because these get made so often
type d map[string]interface{}

//AuthFunc is a function used to determine if the request is authorized
type AuthFunc func(*http.Request) bool

//Options when adding a collection to the admin
type Options struct {
	//Which columns to display/order to display them - nil means all
	Columns []string
}

//Stores info about a specific collection, like the type of the object it
//represents and any options used in specifying the type
type collectionInfo struct {
	Options *Options
	Type    reflect.Type
}

//our main collection mapping variable
var collections = make(map[string]collectionInfo)

//Registers the type/collection pair in the admin. Panics if two types are mapped
//to the same collection
func Register(typ interface{}, coll string, opt *Options) {
	t := reflect.TypeOf(typ)
	if ci, ok := collections[coll]; ok {
		panic(fmt.Sprintf("collection already registered. Had %q->%s . Got %q->%s", coll, ci.Type, coll, t))
	}
	collections[coll] = collectionInfo{opt, t}
}

//Unregisters the information for the colleciton. Panics if you attempt to unregister
//a collection not yet registered.
func Unregister(coll string) {
	if _, ok := collections[coll]; !ok {
		panic(fmt.Sprintf("unregister collection that does not exist: %q", coll))
	}
	delete(collections, coll)
}

//Returns an interface{} boxing a new(T) where T is the type registered
//under the collection name.
func getType(coll string) interface{} {
	t, ok := collections[coll]
	if !ok {
		return nil
	}

	return reflect.New(t.Type).Interface()
}

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

//Admin is an http.Handler for serving up the admin pages
type Admin struct {
	Auth     AuthFunc
	Session  *mgo.Session
	Database string
	Debug    bool

	server *http.ServeMux
}

type Renderer interface {
}

//adminHandler is a type representing a handler function on an *Admin
type adminHandler func(*Admin, http.ResponseWriter, *http.Request)

//routes define the routes for the 
var routes = map[string]adminHandler{
	"/":        (*Admin).Index,
	"/list/":   (*Admin).List,
	"/update/": (*Admin).Update,
	"/create/": (*Admin).Create,
	"/detail/": (*Admin).Detail,
}

//generateMux creates the internal http.ServeMux to dispatch reqeusts to the
//appropriate handler.
func (a *Admin) generateMux() {
	if a.server != nil {
		return
	}

	a.server = http.NewServeMux()
	for r, fn := range routes {
		a.server.Handle(r, a.bind(r, fn))
	}
}

//bind turns an adminHandler into an http.HandlerFunc by closing on the admin
//value on the adminHandler.
func (a *Admin) bind(prefix string, fn adminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, prefix) {
			http.NotFound(w, req)
			return
		}
		req.URL.Path = req.URL.Path[len(prefix):]
		fn(a, w, req)
	}
}

//Returns the mgo.Collection and interface{} which is a *T for sticking data into
func (a *Admin) collFor(coll string) (mgo.Collection, interface{}) {
	return a.Session.DB(a.Database).C(coll), getType(coll)
}

//Presents the detail view for an object in a collection
func (a *Admin) Detail(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id == "" {
		http.NotFound(w, req)
		return
	}

	c, t := a.collFor(coll)

	//load into T
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%v", t)
}

//Presents the index page giving an overall view of the database
func (a *Admin) Index(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll != "" || id != "" {
		http.NotFound(w, req)
		return
	}

}

//Presents a list of objects in a collection matching filtering/sorting criteria
func (a *Admin) List(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id != "" {
		http.NotFound(w, req)
	}

	fmt.Fprintf(w, "%s", coll)
}

//Presents a handler that updates an object and shows the results of the update
func (a *Admin) Update(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id == "" {
		http.NotFound(w, req)
		return
	}

	fmt.Fprintf(w, "%s / %s", coll, id)
}

//Presents a handler that creates an object and shows the results of the create
func (a *Admin) Create(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id != "" {
		http.NotFound(w, req)
	}

	fmt.Fprintf(w, "%s", coll)
}

//ServeHTTP lets *Admin conform to the http.Handler interface for use in web servers
func (a *Admin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if a.Auth != nil && !a.Auth(req) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "You are unauthorized to complete this request")
		return
	}

	//we need that database connection. figure out how to do tests
	if a.Session == nil || a.Database == "" {
		http.Error(w, "Database not set up properly", http.StatusInternalServerError)
		return
	}

	//pass it off to our internal muxer
	a.generateMux()
	a.server.ServeHTTP(w, req)
}

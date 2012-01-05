package admin

import (
	"fmt"
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
func Register(typ interface{}, collection string, options *Options) {
	t := reflect.TypeOf(typ)
	if c, ok := collections[collection]; ok {
		panic(fmt.Sprintf("collection already registered: %s -> %s", c, t))
	}
	collections[collection] = collectionInfo{options, t}
}

//Returns an interface{} that corresponds to a *T where T is the type registered
//under the collection name.
func GetType(coll string) interface{} {
	t, ok := collections[coll]
	if !ok {
		return nil
	}

	return reflect.New(t.Type).Interface()
}

//Admin is an http.Handler for serving up the admin pages
type Admin struct {
	Auth     AuthFunc
	Database interface{} //the mgo connection
}

func (a *Admin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if a.Auth != nil && !a.Auth(req) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "You are unauthorized to complete this request")
		return
	}

	//check the health of the database
	if a.Database == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Database connection broked!")
		return
	}

	r := parseRequest(req)
	//some kind of error.
	if r == nil {
		fmt.Fprintln(w, "Unable to process that request.")
	}

	fmt.Fprintln(w, r)
}

type request struct {
	coll   string                 //the collection
	action action                 //the action
	params map[string]interface{} //the params
}

type action string

var (
	detailAct action = "detail"
	listAct   action = "list"
	updateAct action = "update"
	deleteAct action = "delete"
	createAct action = "create"
)

func parseRequest(req *http.Request) *request {
	chunks := strings.Split(path.Clean(req.URL.Path), "/")

	//this happens when we list them all
	if chunks[0] == "." {
		return &request{action: listAct}
	}

	r := new(request)
	r.coll = chunks[0]
	r.action = action(chunks[1])

	//now we populate params based on the action
	switch r.action {
	case detailAct, deleteAct: //we need an id
		//bail if we dont have an id
		if len(chunks) < 3 {
			return nil
		}
		r.params = d{"id": chunks[2]}
	case updateAct: //we need an id and a form
		//bail if we dont have an id
		if len(chunks) < 3 {
			return nil
		}
		r.params = d{
			"id":   chunks[2],
			"form": req.Form,
		}
	case createAct: //we need a form
		r.params = d{"form": req.Form}
	case listAct: //params for sorting/filtering
		//nothin for now
	default: //this is an unknown action!
		return nil
	}

	return r
}

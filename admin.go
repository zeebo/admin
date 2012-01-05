package admin

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

type AuthFunc func(*http.Request) bool

//Stores info about a specific collection, like the type of the object it
//represents and any options used in specifying the type
type collectionInfo struct {
	Options *Options
	Type    reflect.Type
}

//Options when adding a collection to the admin
type Options struct {
	Columns []string
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
		return
	}

	route := a

	log.Println(req.URL.Path)
}

type route struct {
	coll string
	action 
}



func ()

package admin

import (
	"fmt"
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
	if coll == "" || id == "" {
		http.NotFound(w, req)
		return
	}

	c, t := a.collFor(coll), a.newType(coll)

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

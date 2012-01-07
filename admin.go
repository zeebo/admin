package admin

import (
	"errors"
	"launchpad.net/mgo"
	"net/http"
)

//Admin is an http.Handler for serving up the admin pages
type Admin struct {
	Auth     AuthFunc
	Session  *mgo.Session
	Database string
	Debug    bool
	Renderer Renderer

	//created on demand
	server      *http.ServeMux
	collections map[string]collectionInfo
}

//useful type because these get made so often
type d map[string]interface{}

//AuthFunc is a function used to determine if the request is authorized
type AuthFunc func(*http.Request) bool

//adminHandler is a type representing a handler function on an *Admin
type adminHandler func(*Admin, http.ResponseWriter, *http.Request)

//routes define the routes for the admin
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
		a.server.Handle(r, http.StripPrefix(r, a.bind(fn)))
	}
}

//bind turns an adminHandler into an http.HandlerFunc by closing on the admin
//value on the adminHandler.
func (a *Admin) bind(fn adminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fn(a, w, req)
	}
}

//Returns the mgo.Collection for the specified collection
func (a *Admin) collFor(coll string) mgo.Collection {
	return a.Session.DB(a.Database).C(coll)
}

//ServeHTTP lets *Admin conform to the http.Handler interface for use in web servers
func (a *Admin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if a.Renderer == nil {
		a.Renderer = DefaultRenderer{}
	}

	if a.Auth != nil && !a.Auth(req) {
		a.Renderer.Unauthorized(w, req)
		return
	}

	//ensure a valid database
	if a.Session == nil || a.Database == "" {
		a.Renderer.InternalError(w, req, errors.New("Database not configured"))
		return
	}

	//pass it off to our internal muxer
	a.generateMux()
	a.server.ServeHTTP(w, req)
}

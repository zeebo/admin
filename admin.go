package admin

import (
	"crypto/rand"
	"github.com/zeebo/sign"
	"io"
	"launchpad.net/mgo"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

//useful type because these get made so often
type d map[string]interface{}

//adminHandler is a type representing a handler function on an *Admin.
type adminHandler func(*Admin, http.ResponseWriter, *http.Request)

//Admin is an http.Handler for serving up the admin pages
type Admin struct {
	Auth     Authorizer        //If not nil, admin is auth protected.
	Session  *mgo.Session      //The mongo session for managing.
	Renderer Renderer          //If nil, a default renderer is used to render the admin pages.
	Routes   map[string]string //Routes lets you change the url paths. If nil, uses DefaultRoutes.
	Prefix   string            //The path the admin is mounted to in the handler.
	Key      []byte            //Key for cryptographically signing cookies. Generated if nil.
	Logger   io.Writer         //If nil, os.Stdout is used for logging information.

	//created on demand
	initd       bool
	server      *http.ServeMux
	types       map[string]collectionInfo
	index_cache map[string][]string
	object_id   map[reflect.Type]int
	object_coll map[reflect.Type]string
	auth_cache  map[*http.Request]AuthSession
	logger      *log.Logger
}

//DefaultRoutes is the mapping of actions to url paths.
var DefaultRoutes = map[string]string{
	"index":  "/",
	"list":   "/list/",
	"update": "/update/",
	"create": "/create/",
	"detail": "/detail/",
	"delete": "/delete/",
	"auth":   "/auth/",
}

//routes defines the mapping of type to function for the admin
var routes = map[string]adminHandler{
	"index":  (*Admin).index,
	"list":   (*Admin).list,
	"update": (*Admin).update,
	"create": (*Admin).create,
	"detail": (*Admin).detail,
	"delete": (*Admin).delete,
	"auth":   (*Admin).auth,
}

func (a *Admin) Init() {
	//ensure a valid database
	if a.Session == nil {
		panic("Mongo session not configured")
	}

	//make defaults
	if a.Routes == nil {
		a.Routes = DefaultRoutes
	}
	if a.Key == nil {
		a.generateKey(128) //128 byte key
	}
	if a.Logger == nil {
		a.Logger = os.Stdout
	}
	a.logger = log.New(a.Logger, "ADMIN", log.LstdFlags)

	if a.Renderer == nil {
		a.Renderer = newDefaultRenderer(a.logger)
	}

	required := []string{"index", "list", "update", "create", "detail", "delete", "auth"}
	for _, r := range required {
		if _, ex := a.Routes[r]; !ex {
			panic("Route missing: " + r)
		}
	}

	a.generateMux()
	a.generateIndexCache()

	a.auth_cache = make(map[*http.Request]AuthSession)

	a.initd = true
}

//generateMux creates the internal http.ServeMux to dispatch reqeusts to the
//appropriate handler.
func (a *Admin) generateMux() {
	a.server = http.NewServeMux()
	for key, path := range a.Routes {
		r, fn := path, routes[key]
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

//generateIndexCache generates the values needed for IndexContext and stores
//them for efficient lookup.
func (a *Admin) generateIndexCache() {
	a.index_cache = make(map[string][]string)
	for key := range a.types {
		pieces := strings.Split(key, ".")
		a.index_cache[pieces[0]] = append(a.index_cache[pieces[0]], pieces[1])
	}
}

//generateKey generates a key for cryptographically signing cookie values.
func (a *Admin) generateKey(size int) {
	a.Key = make([]byte, size)
	for p := 0; p < len(a.Key); {
		n, err := rand.Read(a.Key[p:])
		if err != nil {
			panic("Error while generating key: " + err.Error())
		}
		p += n
	}
}

//collFor returns the mgo.Collection for the specified database.collection.
func (a *Admin) collFor(dbcoll string) mgo.Collection {
	pieces := strings.Split(dbcoll, ".")
	return a.Session.DB(pieces[0]).C(pieces[1])
}

//ServeHTTP lets *Admin conform to the http.Handler interface for use in web servers.
func (a *Admin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !a.initd {
		panic("Admin served without calling Init first.")
	}

	//strip off the prefix
	req.URL.Path = req.URL.Path[len(a.Prefix):]

	//if they're going to the auth handler, let them through
	if a.Auth == nil || strings.HasPrefix(req.URL.Path, a.Routes["auth"]) {
		a.server.ServeHTTP(w, req)
		return
	}

	//set up a redirect function to handle adding the redirect cookie
	//and sending them to the login handler
	redirect := func() {
		reverser := Reverser{a}
		http.SetCookie(w, &http.Cookie{
			Name:  "redirect",
			Value: req.URL.Path,
		})
		http.Redirect(w, req, reverser.Login(), http.StatusTemporaryRedirect)
	}

	signer := sign.Signer{a.Key}
	var session AuthSession

	cook, err := req.Cookie("auth")
	if err != nil {
		redirect()
		return
	}

	if err := signer.Unsign(cook.Value, &session, 0); err != nil {
		redirect()
		return
	}

	//store the auth session into our cache
	a.auth_cache[req] = session
	defer delete(a.auth_cache, req)

	a.server.ServeHTTP(w, req)
}

package admin

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//defaultRenderer conforms to the Renderer interface and uses some magic templates
//to create a pretty default interface.
type defaultRenderer struct {
	templates *template.Template
	watcher   chan bool
	mtimes    map[string]time.Time
}

//Lookup returns a template ready to be executed from the template cache, starting
//the watcher goroutine if it has not been started to recompile things at runtime.
func (d *defaultRenderer) Lookup(name string) *template.Template {
	ch := make(chan bool)

	//wait for it to parse
	go d.Watch(ch)
	<-ch

	if d.templates == nil {
		panic("Error compiling templates with nothing to fall back on.")
	}

	t := d.templates.Lookup(name)
	if t == nil {
		panic("Can't find requested template: " + name)
	}

	return t
}

//TemplateDir looks at the environment to find out where the templates live.
//The default value is "./templates"
func (d *defaultRenderer) TemplateDir() string {
	if dir := os.Getenv("ADMIN_TEMPLATE_DIR"); dir != "" {
		return dir
	}
	return "./templates"
}

func (d *defaultRenderer) Watch(parsed chan bool) {
	//only ever spawn one
	select {
	case d.watcher <- true:
	default:
		parsed <- true
		return
	}
	var (
		sentParsed bool
		ticker     = time.NewTicker(1e9) //1 sec
	)
	//do our watching in a forever loop
	for {
		//TODO: detect changes instead of just parsing every time

		t, err := template.ParseGlob(filepath.Join(d.TemplateDir(), "*"))
		if err != nil {
			log.Printf("Error parsing templates: %s", err)
		} else {
			d.templates = t
		}
		if !sentParsed {
			parsed <- true
			sentParsed = true
		}

		<-ticker.C
	}
}

//NotFound presents a basic 404 with no special body.
func (r *defaultRenderer) NotFound(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
	r.Lookup("404").Execute(w, nil)
}

//InternalError presents a basic 500 not suitable for production. Errors should be logged
//and not displayed to the end user.
func (r *defaultRenderer) InternalError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	r.Lookup("internal").Execute(w, err)
}

//Unauthorized presents a simple unauthorized page.
func (r *defaultRenderer) Unauthorized(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	r.Lookup("unauthorized").Execute(w, nil)
}

//Detail presents the detail view of an object.
func (r *defaultRenderer) Detail(w http.ResponseWriter, req *http.Request, c DetailContext) {
	w.Header().Add("Content-Type", "text/html")
	r.Lookup("detail").Execute(w, c)
}

func (r *defaultRenderer) Delete(w http.ResponseWriter, req *http.Request, c DeleteContext) {
	w.Header().Add("Content-Type", "text/html")
	r.Lookup("delete").Execute(w, c)
}

//Index presents an overall view of the database and the managed collections.
func (r *defaultRenderer) Index(w http.ResponseWriter, req *http.Request, c IndexContext) {
	w.Header().Add("Content-Type", "text/html")
	r.Lookup("index").Execute(w, c)
}

//List presents all of the objects of a specific list with the columns and order given by the options
//the type was loaded with.
func (r *defaultRenderer) List(w http.ResponseWriter, req *http.Request, c ListContext) {
	w.Header().Add("Content-Type", "text/html")
	r.Lookup("list").Execute(w, c)
}

//Update presents a success page or the errors when attempting to update an object.
func (r *defaultRenderer) Update(w http.ResponseWriter, req *http.Request, c UpdateContext) {
	w.Header().Add("Content-Type", "text/html")
	r.Lookup("update").Execute(w, c)
}

//Create presents a success page or the errors when attempting to create an object.
func (r *defaultRenderer) Create(w http.ResponseWriter, req *http.Request, c CreateContext) {
	w.Header().Add("Content-Type", "text/html")
	r.Lookup("create").Execute(w, c)
}

package admin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//newDefaultRenderer returns a *defaultRenderer ready to be used.

func newDefaultRenderer(l *log.Logger) *defaultRenderer {
	return &defaultRenderer{
		mtimes:   make(map[string]time.Time),
		newtemp:  make(chan *template.Template),
		currtemp: make(chan *template.Template),
		logger:   l,
	}
}

//defaultRenderer conforms to the Renderer interface and uses some magic templates
//to create a pretty default interface.
type defaultRenderer struct {
	initd    sync.Once
	mtimes   map[string]time.Time
	newtemp  chan *template.Template
	currtemp chan *template.Template
	logger   *log.Logger
}

//init is called once on a defaultRenderer. Sets up the system for watching
//the directory of templates.
func (d *defaultRenderer) init() {
	d.initd.Do(func() {
		//seed the parsing
		tmpl, err := d.parse()
		if err != nil {
			panic(err)
		}

		//setup watchers
		go d.watch()
		go d.sender(tmpl)
	})
}

//sender is a simple function to always send out the most current template (with
//very low probability that it stays outdated for long.)
func (d *defaultRenderer) sender(curr *template.Template) {
	for {
		select {
		case curr = <-d.newtemp:
		case d.currtemp <- curr:
		}
	}
}

//Lookup returns a template ready to be executed from the template cache, starting
//the watcher goroutine if it has not been started to recompile things at runtime.
func (d *defaultRenderer) Lookup(name string) *template.Template {
	d.init()

	t := (<-d.currtemp).Lookup(name)
	if t == nil {
		panic(fmt.Errorf("Can't find requested template: %s", name))
	}

	return t
}

//dir looks at the environment to find out where the templates live.
//The default value is "./templates"
func (d *defaultRenderer) dir() string {
	if dir := os.Getenv("ADMIN_TEMPLATE_DIR"); dir != "" {
		return dir
	}
	return "templates"
}

//updateMtimes globs the template directory for files and checks their modified
//times to see if they differ, updating the mtimes cache.
func (d *defaultRenderer) updateMtimes() (bool, error) {
	files, err := filepath.Glob(filepath.Join(d.dir(), "*"))
	if err != nil {
		return false, err
	}

	var changes bool
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return false, err
		}

		mtime := info.ModTime()
		if pmtime, ex := d.mtimes[file]; !ex || mtime != pmtime {
			changes = true
			d.mtimes[file] = mtime
		}
	}

	return changes, nil
}

//watch watches the template directory every second for modifications and
//recompiles the templates.
func (d *defaultRenderer) watch() {
	var ticker = time.NewTicker(1e9) //1 sec
	//do our watching in a forever loop
	for {
		<-ticker.C

		tmpl, err := d.parse()
		if err != nil {
			d.logger.Printf("Error parsing templates: %s", err)
			continue
		}

		if tmpl != nil {
			d.newtemp <- tmpl
			d.logger.Print("Templates updated.")
		}
	}
}

//parse checks the modified times and parses the template directory if required.
func (d *defaultRenderer) parse() (tmpl *template.Template, err error) {
	changed, err := d.updateMtimes()
	if err != nil {
		return
	}
	if !changed {
		return
	}

	tmpl = template.New("base")
	tmpl.Funcs(template.FuncMap{
		"noescape": func(a ...interface{}) string {
			return fmt.Sprint(a...)
		},
	})

	tmpl, err = tmpl.ParseGlob(filepath.Join(d.dir(), "*"))
	return
}

//NotFound presents a basic 404 with no special body.
func (r *defaultRenderer) NotFound(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
	if err := r.Lookup("404").Execute(w, nil); err != nil {
		panic(err)
	}
}

//InternalError presents a basic 500 not suitable for production. Errors should be logged
//and not displayed to the end user.
func (r *defaultRenderer) InternalError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	if err := r.Lookup("internal").Execute(w, err); err != nil {
		panic(err)
	}
}

//Unauthorized presents a simple unauthorized page.
func (r *defaultRenderer) Unauthorized(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	if err := r.Lookup("unauthorized").Execute(w, nil); err != nil {
		panic(err)
	}
}

//Detail presents the detail view of an object.
func (r *defaultRenderer) Detail(w http.ResponseWriter, req *http.Request, c DetailContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("detail").Execute(w, c); err != nil {
		panic(err)
	}
}

func (r *defaultRenderer) Delete(w http.ResponseWriter, req *http.Request, c DeleteContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("delete").Execute(w, c); err != nil {
		panic(err)
	}
}

//Index presents an overall view of the database and the managed collections.
func (r *defaultRenderer) Index(w http.ResponseWriter, req *http.Request, c BaseContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("index").Execute(w, c); err != nil {
		panic(err)
	}
}

//List presents all of the objects of a specific list with the columns and order given by the options
//the type was loaded with.
func (r *defaultRenderer) List(w http.ResponseWriter, req *http.Request, c ListContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("list").Execute(w, c); err != nil {
		panic(err)
	}
}

//Update presents a success page or the errors when attempting to update an object.
func (r *defaultRenderer) Update(w http.ResponseWriter, req *http.Request, c UpdateContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("update").Execute(w, c); err != nil {
		panic(err)
	}
}

//Create presents a success page or the errors when attempting to create an object.
func (r *defaultRenderer) Create(w http.ResponseWriter, req *http.Request, c CreateContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("create").Execute(w, c); err != nil {
		panic(err)
	}
}

//Authorize presents a page with a login form for authorizing a user.
func (r *defaultRenderer) Authorize(w http.ResponseWriter, req *http.Request, c AuthorizeContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("login").Execute(w, c); err != nil {
		panic(err)
	}
}

//LoggedOut presents a page thanking the user for spending time with the site.
func (r *defaultRenderer) LoggedOut(w http.ResponseWriter, req *http.Request, c BaseContext) {
	w.Header().Add("Content-Type", "text/html")
	if err := r.Lookup("logout").Execute(w, c); err != nil {
		panic(err)
	}
}

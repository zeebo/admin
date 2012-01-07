package admin

import (
	"flag"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"log"
	"net/http"
	"os/exec"
	"testing"
)

var (
	database = flag.String("db", "admin_test", "Database for mongo connnection")
	reload   = flag.Bool("load", false, "Run mongoimport on the json file for the database")
	jsonfile = flag.String("json", "admin_test.json", "Json file for loading into the database")
)

func init() {
	flag.Parse()

	//Import: mongoimport --drop -d admin_test -c T admin_test.json
	//Export: mongoexport -d admin_test -c T > admin_test.json

	//before commit:
	//mongoexport -d admin_test -c T > admin_test.json
	//go test -load
	//git commit -a -m 'msg'

	if *reload {
		cmd := exec.Command("mongoimport", "--drop", "-d", *database, "-c", "T", *jsonfile)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

type T struct {
	ID bson.ObjectId `bson:"_id"`
}

func TestRegister(t *testing.T) {
	h := &Admin{}

	h.Register(T{}, "T", nil)

	ans := h.newType("T")
	if _, ok := ans.(*T); !ok {
		t.Fatalf("Type incorrect. Expected *admin.T, got %T", ans)
	}
}

func TestRegisterFail(t *testing.T) {
	h := &Admin{}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("No panic when attempting to reregister type")
		}
	}()

	h.Register(T{}, "T", nil)
	h.Register(T{}, "T", nil)
}

func TestNewTypeNewInstance(t *testing.T) {
	h := &Admin{}

	h.Register(T{}, "T", nil)

	ans1 := h.newType("T")
	ans2 := h.newType("T")
	if ans1.(*T) == ans2.(*T) {
		t.Fatal("getType returned identical instances")
	}
}

func TestUnauthorized(t *testing.T) {
	h := &Admin{
		Auth: func(*http.Request) bool { return false },
	}

	w := Get(t, h, "/")
	if w.Status != http.StatusUnauthorized {
		t.Fatalf("Failed being unauthorized. Got %d", w.Status)
	}
}

func TestAuthorized(t *testing.T) {
	h := &Admin{
		Auth: func(*http.Request) bool { return true },
	}

	w := Get(t, h, "/")
	if w.Status == http.StatusUnauthorized {
		t.Fatalf("Failed being authorized. Got %d", w.Status)
	}
}

func TestInvalidDetail(t *testing.T) {
	s, _ := mgo.Mongo("")
	h := &Admin{
		Database: *database,
		Session:  s,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/detail/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Detail did not 404 without collection. Got %d", w.Status)
	}

	w = Get(t, h, "/detail/T/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Detail did not 404 without id. Got %d", w.Status)
	}
}

func TestInvalidIndex(t *testing.T) {
	s, _ := mgo.Mongo("")
	h := &Admin{
		Database: *database,
		Session:  s,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/foobar")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Index did not 404. Got %d", w.Status)
	}
}

func TestInvalidList(t *testing.T) {
	s, _ := mgo.Mongo("")
	h := &Admin{
		Database: *database,
		Session:  s,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/list/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("List did not 404 without collection. Got %d", w.Status)
	}
}

func TestInvalidUpdate(t *testing.T) {
	s, _ := mgo.Mongo("")
	h := &Admin{
		Database: *database,
		Session:  s,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/update/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Update did not 404 without collection. Got %d", w.Status)
	}

	w = Get(t, h, "/update/T/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Update did not 404 without id. Got %d", w.Status)
	}
}

func TestInvalidCreate(t *testing.T) {
	s, _ := mgo.Mongo("")
	h := &Admin{
		Database: *database,
		Session:  s,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/create/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Create did not 404 without collection. Got %d", w.Status)
	}
}

func TestCorrectRenderIndex(t *testing.T) {
	s, _ := mgo.Mongo("")
	r := &TestRenderer{}
	h := &Admin{
		Database: *database,
		Session:  s,
		Renderer: r,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on Index. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Index" {
		t.Fatalf("Wrong Renderer called. Expected Index got %s", r.Last().Type)
	}
}

func TestCorrectRenderList(t *testing.T) {
	s, err := mgo.Mongo("localhost")
	if err != nil {
		t.Fatalf("Unable to connect to test database: %s", err)
	}
	r := &TestRenderer{}
	h := &Admin{
		Database: *database,
		Session:  s,
		Renderer: r,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/list/T/")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on List. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "List" {
		t.Fatalf("Wrong Renderer called. Expected List got %s", r.Last().Type)
	}
}

func TestCorrectRenderUpdate(t *testing.T) {
	s, err := mgo.Mongo("localhost")
	if err != nil {
		t.Fatalf("Unable to connect to test database: %s", err)
	}
	r := &TestRenderer{}
	h := &Admin{
		Database: *database,
		Session:  s,
		Renderer: r,
	}
	h.Register(T{}, "T", nil)
	var w *TestResponseWriter

	w = Get(t, h, "/update/T/4f07c34779bf562daff8640c")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on Update. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Update" {
		t.Fatalf("Wrong Renderer called. Expected Update got %s", r.Last().Type)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Error checking for correct object. %s", err)
		}
	}()

	if id := r.Last().Params.(UpdateContext).Object.(*T).ID.Hex(); id != "4f07c34779bf562daff8640c" {
		t.Fatalf("Update returned the wrong object. Expected 4f07c34779bf562daff8640c got %s", id)
	}
}

func TestCorrectRenderCreate(t *testing.T) {
	s, err := mgo.Mongo("localhost")
	if err != nil {
		t.Fatalf("Unable to connect to test database: %s", err)
	}
	r := &TestRenderer{}
	h := &Admin{
		Database: *database,
		Session:  s,
		Renderer: r,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/create/T/")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on Create. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Create" {
		t.Fatalf("Wrong Renderer called. Expected Create got %s", r.Last().Type)
	}
}

func TestCorrectRenderDetail(t *testing.T) {
	s, err := mgo.Mongo("localhost")
	if err != nil {
		t.Fatalf("Unable to connect to test database: %s", err)
	}
	r := &TestRenderer{}
	h := &Admin{
		Database: *database,
		Session:  s,
		Renderer: r,
	}
	h.Register(T{}, "T", nil)
	var w *TestResponseWriter

	w = Get(t, h, "/detail/T/4f07c34779bf562daff8640c")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on Detail. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Detail" {
		t.Fatalf("Wrong Renderer called. Expected Detail got %s", r.Last().Type)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Error checking for correct object. %s", err)
		}
	}()

	if id := r.Last().Params.(DetailContext).Object.(*T).ID.Hex(); id != "4f07c34779bf562daff8640c" {
		t.Fatalf("Detail returned the wrong object. Expected 4f07c34779bf562daff8640c got %s", id)
	}
}

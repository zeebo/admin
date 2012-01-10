package admin

import (
	"fmt"
	"launchpad.net/gobson/bson"
	"net/http"
	"testing"
)

type T struct {
	ID bson.ObjectId `bson:"_id"`
}

func (t T) GetTemplate() string {
	return ``
}

type T2 struct {
	ID bson.ObjectId `bson:"_id"`
	V  int           `bson:"v"`
}

func (t T2) GetTemplate() string {
	return ``
}

type T3 struct {
}

func (t T3) GetTemplate() string {
	return `{{`
}

func TestRegisterWorks(t *testing.T) {
	h := &Admin{}

	h.Register(T{}, "admin_test.T", nil)

	ans := h.newType("admin_test.T")
	if _, ok := ans.(*T); !ok {
		t.Fatalf("Type incorrect. Expected *admin.T, got %T", ans)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	h := &Admin{}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("No panic when attempting to reregister type")
		}
	}()

	h.Register(T{}, "admin_test.T", nil)
	h.Register(T{}, "admin_test.T", nil)
}

func TestRegisterBadTemplate(t *testing.T) {
	h := &Admin{}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("No panic when attempting to register bad template")
		}
	}()

	h.Register(T3{}, "admin_test.T3", nil)
}

func TestNewTypeNewInstance(t *testing.T) {
	h := &Admin{}

	h.Register(T{}, "admin_test.T", nil)

	ans1 := h.newType("admin_test.T")
	ans2 := h.newType("admin_test.T")
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

func TestDetailInvalid(t *testing.T) {
	h := &Admin{
		Session: session,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/detail/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Detail did not 404 without collection. Got %d", w.Status)
	}

	w = Get(t, h, "/detail/admin_test.T/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Detail did not 404 without id. Got %d", w.Status)
	}

	w = Get(t, h, "/detail/admin_test.T/ffffffffffffffffffffffff")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Wrong return type on Update. Expected 200 got %d", w.Status)
	}
}

func TestIndexInvalid(t *testing.T) {
	h := &Admin{
		Session: session,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/foobar")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Index did not 404. Got %d", w.Status)
	}
}

func TestListInvalid(t *testing.T) {
	h := &Admin{
		Session: session,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/list/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("List did not 404 without collection. Got %d", w.Status)
	}
}

func TestUpdateInvalid(t *testing.T) {
	h := &Admin{
		Session: session,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/update/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Update did not 404 without collection. Got %d", w.Status)
	}

	w = Get(t, h, "/update/admin_test.T/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Update did not 404 without id. Got %d", w.Status)
	}

	w = Get(t, h, "/update/admin_test.T/ffffffffffffffffffffffff")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Wrong return type on Update. Expected 200 got %d", w.Status)
	}
}

func TestCreateInvalid(t *testing.T) {
	h := &Admin{
		Session: session,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/create/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Create did not 404 without collection. Got %d", w.Status)
	}
}

func TestIndexCorrectRender(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
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

func TestListCorrectRender(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T{}, "admin_test.T", nil)
	var w *TestResponseWriter

	w = Get(t, h, "/list/admin_test.T/")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on List. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "List" {
		t.Fatalf("Wrong Renderer called. Expected List got %s", r.Last().Type)
	}
}

func TestUpdateCorrectRender(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T{}, "admin_test.T", nil)
	var w *TestResponseWriter

	w = Get(t, h, "/update/admin_test.T/4f07c34779bf562daff8640c")
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

func TestCreateCorrectRender(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T{}, "admin_test.T", nil)
	var w *TestResponseWriter

	w = Get(t, h, "/create/admin_test.T/")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on Create. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Create" {
		t.Fatalf("Wrong Renderer called. Expected Create got %s", r.Last().Type)
	}
}

func TestDetailCorrectRender(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T{}, "admin_test.T", nil)
	var w *TestResponseWriter

	w = Get(t, h, "/detail/admin_test.T/4f07c34779bf562daff8640c")
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

func TestRegisterNoDatabase(t *testing.T) {
	h := &Admin{}
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("No panic with invalid type")
		}
	}()
	h.Register(T{}, "T", nil)
}

func TestUpdateUnknownCollection(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/update/unknown.T/4f07c34779bf562daff8640c")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Expected 404 got %d", w.Status)
	}
	if r.Last().Type != "NotFound" {
		t.Fatalf("Wrong Renderer called. Expected NotFound got %s", r.Last().Type)
	}
}

func TestDetailUnknownCollection(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/detail/unknown.T/4f07c34779bf562daff8640c")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Expected 404 got %d", w.Status)
	}
	if r.Last().Type != "NotFound" {
		t.Fatalf("Wrong Renderer called. Expected NotFound got %s", r.Last().Type)
	}
}

func TestListUnknownCollection(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/list/unknown.T/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Expected 404 got %d", w.Status)
	}
	if r.Last().Type != "NotFound" {
		t.Fatalf("Wrong Renderer called. Expected NotFound got %s", r.Last().Type)
	}
}

func TestCreateUnknownCollection(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/create/unknown.T/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Expected 404 got %d", w.Status)
	}
	if r.Last().Type != "NotFound" {
		t.Fatalf("Wrong Renderer called. Expected NotFound got %s", r.Last().Type)
	}
}

func TestListReturns(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T2{}, "admin_test.T2", nil)

	Get(t, h, "/list/admin_test.T2/")
	context := r.Last().Params.(ListContext)
	for i, obj := range context.Objects {
		if obj.(*T2).V != i {
			t.Fatalf("Expected object %d. Got %d", i, obj.(*T2).V)
		}
	}
}

func TestListNumPage(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T2{}, "admin_test.T2", nil)

	test := func(page, numpage int) {
		Get(t, h, fmt.Sprintf("/list/admin_test.T2/?page=%d&numpage=%d", page, numpage))
		context := r.Last().Params.(ListContext)
		for i, obj := range context.Objects {
			if n := i + ((page - 1) * numpage); obj.(*T2).V != n {
				t.Fatalf("Expected object %d. Got %d", n, obj.(*T2).V)
			}
		}
		if len(context.Objects) != numpage {
			t.Fatalf("Expected %d objects on page. Got %d", numpage, len(context.Objects))
		}
	}

	for i := 1; i < 50; i++ {
		for j := 1; i*j < 50; j++ {
			test(i, j)
		}
	}
}

func TestListInvalidParams(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T2{}, "admin_test.T2", nil)

	//default to page 1
	Get(t, h, "/list/admin_test.T2/?page=-1")
	context := r.Last().Params.(ListContext)
	for i, obj := range context.Objects {
		if obj.(*T2).V != i {
			t.Fatalf("Expected object %d. Got %d", i, obj.(*T2).V)
		}
	}

	//default to 20 items
	Get(t, h, "/list/admin_test.T2/?page=-1")
	context = r.Last().Params.(ListContext)
	if len(context.Objects) != 20 {
		t.Fatalf("Expected 20 objects on page. Got %d", len(context.Objects))
	}
}

func TestListSortingOrder(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T2{}, "admin_test.T2", nil)

	Get(t, h, "/list/admin_test.T2/?sort_v=desc")
	context := r.Last().Params.(ListContext)
	for i, obj := range context.Objects {
		if obj.(*T2).V != 49-i {
			t.Fatalf("Expected object %d. Got %d", 49-i, obj.(*T2).V)
		}
	}

	Get(t, h, "/list/admin_test.T2/?sort_v=desc&page=2")
	context = r.Last().Params.(ListContext)
	for i, obj := range context.Objects {
		if obj.(*T2).V != 29-i {
			t.Fatalf("Expected object %d. Got %d", 49-i, obj.(*T2).V)
		}
	}

	Get(t, h, "/list/admin_test.T2/?sort_v=desc&page=2&numpage=5")
	context = r.Last().Params.(ListContext)
	for i, obj := range context.Objects {
		if obj.(*T2).V != 44-i {
			t.Fatalf("Expected object %d. Got %d", 49-i, obj.(*T2).V)
		}
	}

	Get(t, h, "/list/admin_test.T2/?sort_v=asc")
	context = r.Last().Params.(ListContext)
	for i, obj := range context.Objects {
		if obj.(*T2).V != i {
			t.Fatalf("Expected object %d. Got %d", 49-i, obj.(*T2).V)
		}
	}
}

func TestListSortingInvalid(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T2{}, "admin_test.T2", nil)

	Get(t, h, "/list/admin_test.T2/?sort_no_field=desc")
	context := r.Last().Params.(ListContext)
	for i, obj := range context.Objects {
		if obj.(*T2).V != i {
			t.Fatalf("Expected object %d. Got %d", 49-i, obj.(*T2).V)
		}
	}
}

package admin

import (
	"launchpad.net/mgo"
	"net/http"
	"testing"
)

type T struct{}

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
		Database: "admin_test",
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
		Database: "admin_test",
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
		Database: "admin_test",
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
		Database: "admin_test",
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
		Database: "admin_test",
		Session:  s,
	}
	var w *TestResponseWriter

	w = Get(t, h, "/create/")
	if w.Status != http.StatusNotFound {
		t.Fatalf("Create did not 404 without collection. Got %d", w.Status)
	}
}

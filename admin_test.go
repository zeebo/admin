package admin

import (
	"net/http"
	"testing"
)

func TestAdminRegisterNoID(t *testing.T) {
	h := &Admin{
		Renderer: &TestRenderer{},
	}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("No panic when attempting to register a type with no ID")
		}
	}()

	h.Register(T7{}, "admin_test.T7", nil)
}

func TestAdminRegisterWorks(t *testing.T) {
	h := &Admin{
		Renderer: &TestRenderer{},
	}

	h.Register(T{}, "admin_test.T", nil)

	ans := h.newType("admin_test.T")
	if _, ok := ans.(*T); !ok {
		t.Fatalf("Type incorrect. Expected *admin.T, got %T", ans)
	}
}

func TestAdminRegisterPointer(t *testing.T) {
	h := &Admin{
		Renderer: &TestRenderer{},
	}

	h.Register(&T{}, "admin_test.T", nil)

	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	h.newType("admin_test.T")
}

func TestAdminRegisterDuplicate(t *testing.T) {
	h := &Admin{
		Renderer: &TestRenderer{},
	}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("No panic when attempting to reregister type")
		}
	}()

	h.Register(T{}, "admin_test.T", nil)
	h.Register(T{}, "admin_test.T", nil)
}

func TestAdminRegisterInvalidType(t *testing.T) {
	h := &Admin{
		Renderer: &TestRenderer{},
	}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("No panic when attempting to register an unloadable type")
		}
	}()

	h.Register(T4{}, "admin_test.T4", nil)
}

func TestAdminRegisterCustomLoader(t *testing.T) {
	h := &Admin{
		Renderer: &TestRenderer{},
	}

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Error registering a good type:", err)
		}
	}()

	h.Register(T5{}, "admin_test.T5", nil)

	ans := h.newType("admin_test.T5")
	if _, ok := ans.(*T5); !ok {
		t.Fatalf("Type incorrect. Expected *admin.T, got %T", ans)
	}
}

func TestAdminNewTypeNewInstance(t *testing.T) {
	h := &Admin{
		Renderer: &TestRenderer{},
	}

	h.Register(T{}, "admin_test.T", nil)

	ans1 := h.newType("admin_test.T")
	ans2 := h.newType("admin_test.T")
	if ans1.(*T) == ans2.(*T) {
		t.Fatal("getType returned identical 1nstances")
	}
}

func TestAdminRegisterNoDatabase(t *testing.T) {
	h := &Admin{}
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("No panic with invalid type")
		}
	}()
	h.Register(T{}, "T", nil)
}

func TestAdminCustomPaths(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Routes: map[string]string{
			"index":  "/1/",
			"list":   "/2/",
			"update": "/3/",
			"create": "/4/",
			"detail": "/5/",
			"delete": "/6/",
			"auth":   "/7/",
		},
		Session:  session,
		Renderer: r,
	}
	h.Register(T{}, "admin_test.T", nil)
	var w *TestResponseWriter

	w = Get(t, h, "/1/")
	if w.Status != http.StatusOK {
		t.Errorf("Wrong return type on Index. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Index" {
		t.Errorf("Wrong Renderer called. Expected Index got %s", r.Last().Type)
	}

	w = Get(t, h, "/2/admin_test.T/")
	if w.Status != http.StatusOK {
		t.Errorf("Wrong return type on List. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "List" {
		t.Errorf("Wrong Renderer called. Expected List got %s", r.Last().Type)
	}

	w = Get(t, h, "/3/admin_test.T/4f07c34779bf562daff8640c")
	if w.Status != http.StatusOK {
		t.Errorf("Wrong return type on Update. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Update" {
		t.Errorf("Wrong Renderer called. Expected Update got %s", r.Last().Type)
	}

	w = Get(t, h, "/4/admin_test.T/")
	if w.Status != http.StatusOK {
		t.Errorf("Wrong return type on Create. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Create" {
		t.Errorf("Wrong Renderer called. Expected Create got %s", r.Last().Type)
	}

	w = Get(t, h, "/5/admin_test.T/4f07c34779bf562daff8640c")
	if w.Status != http.StatusOK {
		t.Errorf("Wrong return type on Detail. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Detail" {
		t.Errorf("Wrong Renderer called. Expected Detail got %s", r.Last().Type)
	}

	w = Get(t, h, "/6/admin_test.T/4f07c34779bf562daff8640c")
	if w.Status != http.StatusOK {
		t.Errorf("Wrong return type on Delete. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Delete" {
		t.Errorf("Wrong Renderer called. Expected Delete got %s", r.Last().Type)
	}

	w = Get(t, h, "/7/login")
	if w.Status != http.StatusOK {
		t.Errorf("Wrong return type on Auth. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Authorize" {
		t.Errorf("Wrong Renderer called. Expected Authorize got %s", r.Last().Type)
	}

}

func TestAdminMissingRoutes(t *testing.T) {
	h := &Admin{
		Routes:   map[string]string{},
		Session:  session,
		Renderer: &TestRenderer{},
	}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("Expected error.")
		}
	}()

	Get(t, h, "/foo")
}

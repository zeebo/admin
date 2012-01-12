package admin

import (
	"fmt"
	"launchpad.net/gobson/bson"
	"net/http"
	"net/url"
	"testing"
)

func TestAdminPostCreate(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T6{}, "admin_test.T6", nil)
	var w *TestResponseWriter

	w = Post(t, h, "/create/admin_test.T6/", url.Values{
		"X": {"20"},
		"Y": {"foo"},
		"Z": {"true"},
	})
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on create. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Create" {
		t.Fatalf("Wrong Renderer called. Expected Create got %s", r.Last().Type)
	}

	//lets check all the context values
	params := r.Last().Params.(CreateContext)
	if !params.Success {
		t.Fatalf("Unsucessful create.\n%s", params)
	}
	if !params.Attempted {
		t.Fatalf("Attempted false. Expected true.\n%s", params)
	}

	values := params.Form.context.Values

	defer session.DB("admin_test").C("T6").Remove(d{"_id": bson.ObjectIdHex(values["ID"])})

	if values["X"] != "20" {
		t.Fatalf("X: Expected %q. Got %q.", "20", values["X"])
	}

	if values["Y"] != "foo" {
		t.Fatalf("Y: Expected %q. Got %q.", "foo", values["X"])
	}

	if values["Z"] != "true" {
		t.Fatalf("Z: Expected %q. Got %q.", "true", values["X"])
	}

	//lets ask detail about values["ID"]
	w = Get(t, h, fmt.Sprintf("/detail/admin_test.T6/%s", values["ID"]))
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on Detail. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Detail" {
		t.Fatalf("Wrong Renderer called. Expected Detail got %s", r.Last().Type)
	}

	//cowboy type assertions
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Error checking conditions: %s", err)
		}
	}()

	//compare values in the object to known values
	obj := r.Last().Params.(DetailContext).Object.(*T6)
	if obj.X != 20 {
		t.Fatalf("X: Expected %d. Got %d.", 20, obj.X)
	}
	if obj.Y != "foo" {
		t.Fatalf("Y: Expected %s. Got %s.", "foo", obj.X)
	}
	if obj.Z != true {
		t.Fatalf("Z: Expected %v. Got %v.", true, obj.X)
	}

	//delete the item from the database
}

func TestAdminPostUpdate(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T6{}, "admin_test.T6", nil)
	var w *TestResponseWriter

	//revert to original after
	defer session.DB("admin_test").C("T6").Update(d{"_id": bson.ObjectIdHex("4f0ee3600888a1b6646199bd")}, d{"x": 20, "y": "foo", "z": true})

	w = Post(t, h, "/update/admin_test.T6/4f0ee3600888a1b6646199bd", url.Values{
		"X": {"30"},
		"Y": {"foob"},
		"Z": {"false"},
	})
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on create. expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Update" {
		t.Fatalf("Wrong Renderer called. Expected Update got %s", r.Last().Type)
	}

	w = Get(t, h, "/detail/admin_test.T6/4f0ee3600888a1b6646199bd")
	if w.Status != http.StatusOK {
		t.Fatalf("Wrong return type on Detail. Expected 200 got %d", w.Status)
	}
	if r.Last().Type != "Detail" {
		t.Fatalf("Wrong Renderer called. Expected Detail got %s", r.Last().Type)
	}

	//cowboy type assertions
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Error checking conditions: %s", err)
		}
	}()

	//compare values in the object to known values
	obj := r.Last().Params.(DetailContext).Object.(*T6)
	if obj.X != 30 {
		t.Fatalf("X: Expected %d. Got %d.", 30, obj.X)
	}
	if obj.Y != "foob" {
		t.Fatalf("Y: Expected %s. Got %s.", "foob", obj.X)
	}
	if obj.Z != false {
		t.Fatalf("Z: Expected %v. Got %v.", false, obj.X)
	}
}

func TestAdminEveryAction(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T6{}, "admin_test.T6", nil)

	//cowboy type assertions
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Error while attempting to do everything: %s", err)
		}
	}()

	c := session.DB("admin_test").C("T6")

	Post(t, h, "/create/admin_test.T6/", url.Values{
		"X": {"2"},
		"Y": {"new"},
		"Z": {"false"},
	})
	id := r.Last().Params.(CreateContext).Form.context.Values["ID"]

	//lets check it exists in the database
	if n, err := c.Find(d{"_id": bson.ObjectIdHex(id)}).Count(); n != 1 || err != nil {
		t.Fatalf("Expected %d objects in the database. Got %d.\nError: %v", 1, n, err)
	}

	//make sure we can get it out from the web
	Get(t, h, fmt.Sprintf("/detail/admin_test.T6/%s", id))
	obj1 := r.Last().Params.(DetailContext).Object.(*T6)
	if obj1.X != 2 || obj1.Y != "new" || obj1.Z != false {
		t.Fatalf("Expected %v\nGot %v\n", T6{"", 2, "new", false}, obj1)
	}

	//time to update
	Post(t, h, fmt.Sprintf("/update/admin_test.T6/%s", id), url.Values{
		"X": {"20"},
		"Y": {"newt"},
		"Z": {"true"},
	})
	obj2 := r.Last().Params.(UpdateContext).Object.(*T6)
	if obj2.X != 20 || obj2.Y != "newt" || obj2.Z != true {
		t.Fatalf("Expected %v\nGot %v\n", T6{"", 20, "newt", true}, obj1)
	}

	//now lets delete it
	Get(t, h, fmt.Sprintf("/delete/admin_test.T6/%s?_sure=yes", id))
	if n, err := c.Find(d{"_id": bson.ObjectIdHex(id)}).Count(); n != 0 || err != nil {
		t.Fatalf("Expected %d objects in the database. Got %d.\nError: %v", 0, n, err)
	}
}

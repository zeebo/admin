package admin

import (
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

	id := params.Form.context.Values["ID"]
	t.Log(id)
	//now lets request it back out from detail

	//TODO: create a way to get the ID of the object just inserted. Maybe
	//ask on mgo mailing list?

	session.DB("admin_test").C("T6").DropCollection()
}

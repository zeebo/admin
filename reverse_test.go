package admin

import (
	"launchpad.net/gobson/bson"
	"testing"
)

func TestReverseIdFor(t *testing.T) {
	h := &Admin{}
	r := Reverser{h}
	h.Register(T{}, "admin_test.T", nil)

	var x T = T{bson.ObjectIdHex("ffffffffffffffffffffffff")}

	if r.idFor(x) != "ffffffffffffffffffffffff" {
		t.Fatalf("Expected %q. Got %q.", "ffffffffffffffffffffffff", r.idFor(x))
	}

	if r.idFor(&x) != "ffffffffffffffffffffffff" {
		t.Fatalf("Expected %q. Got %q.", "ffffffffffffffffffffffff", r.idFor(x))
	}
}

func TestReverseCollFor(t *testing.T) {
	h := &Admin{}
	r := Reverser{h}
	h.Register(T{}, "admin_test.T", nil)

	if c := r.collFor(T{}); c != "admin_test.T" {
		t.Fatalf("Expected %q. Got %q.", "admin_test.T", c)
	}
}

func TestReversePrefix(t *testing.T) {
	h := &Admin{
		Prefix: "/admin",
	}
	r := Reverser{h}
	h.Register(T{}, "admin_test.T", nil)

	if r.Index() != "/admin" {
		t.Fatalf("Expected %q. Got %q.", "/admin/", r.Index())
	}
}

func TestReverseDefaultObj(t *testing.T) {
	h := &Admin{}
	r := Reverser{h}
	h.Register(T{}, "admin_test.T", nil)
	h.generateMux()
	var x T = T{bson.ObjectIdHex("ffffffffffffffffffffffff")}

	if c, e := r.CreateObj(x), "/create/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.DetailObj(x), "/detail/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.ListObj(x), "/list/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.DeleteObj(x), "/delete/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.UpdateObj(x), "/update/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
}

func TestReverseCustomObj(t *testing.T) {
	h := &Admin{
		Routes: map[string]string{
			"create": "/1/",
			"detail": "/2/",
			"list":   "/3/",
			"delete": "/4/",
			"update": "/5/",
			"auth":   "/6/",
			"index":  "/",
		},
	}
	r := Reverser{h}
	h.Register(T{}, "admin_test.T", nil)
	h.generateMux()
	var x T = T{bson.ObjectIdHex("ffffffffffffffffffffffff")}

	if c, e := r.CreateObj(x), "/1/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.DetailObj(x), "/2/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.ListObj(x), "/3/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.DeleteObj(x), "/4/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.UpdateObj(x), "/5/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
}

func TestReverseDefaultSpecified(t *testing.T) {
	h := &Admin{}
	r := Reverser{h}
	coll, id := "admin_test.T", "ffffffffffffffffffffffff"

	h.Register(T{}, coll, nil)
	h.generateMux()

	if c, e := r.Create(coll), "/create/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.Detail(coll, id), "/detail/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.List(coll), "/list/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.Delete(coll, id), "/delete/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.Update(coll, id), "/update/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
}

func TestReverseCustomSpecified(t *testing.T) {
	h := &Admin{
		Routes: map[string]string{
			"create": "/1/",
			"detail": "/2/",
			"list":   "/3/",
			"delete": "/4/",
			"update": "/5/",
			"auth":   "/6/",
			"index":  "/",
		},
	}
	r := Reverser{h}
	coll, id := "admin_test.T", "ffffffffffffffffffffffff"
	h.Register(T{}, coll, nil)
	h.generateMux()

	if c, e := r.Create(coll), "/1/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.Detail(coll, id), "/2/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.List(coll), "/3/admin_test.T"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.Delete(coll, id), "/4/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
	if c, e := r.Update(coll, id), "/5/admin_test.T/ffffffffffffffffffffffff"; c != e {
		t.Fatalf("Expected %q. Got %q.", e, c)
	}
}

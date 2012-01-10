package admin

import (
	"launchpad.net/mgo"
	"log"
	"net/http"
)

func ExampleAdmin_Register() {
	a := &Admin{
		Session: session,
	}
	a.Register(T{}, "my_database.T", nil)
}

func ExampleAdmin_ServeHTTP() {
	a := &Admin{
		Session: session,
	}
	a.Register(T{}, "my_database.T", nil)

	http.Handle("/admin/", a)
	if err := http.ListenAndServe(":11223", nil); err != nil {
		log.Fatal(err)
	}
}

func ExampleAdmin_Unregister() {
	a := &Admin{
		Session: session,
	}
	a.Register(T{}, "my_database.T", nil)
	a.Unregister("my_database.T")
}

func ExampleAdmin() {
	session, err := mgo.Mongo("mongo://my_mongo_server")
	if err != nil {
		log.Fatal(err)
	}

	a := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
		Auth:     func(req *http.Request) bool { return true },
		Routes: map[string]string{
			"index":  "/1/",
			"list":   "/2/",
			"update": "/3/",
			"create": "/4/",
			"detail": "/5/",
		},
	}

	a.Register(T{}, "database.collection", &Options{
		Columns: []string{"first", "second", "fifth"},
	})

	http.Handle("/admin/", http.StripPrefix("/admin/", a))
	if err := http.ListenAndServe(":11223", nil); err != nil {
		log.Fatal(err)
	}
}

func ExampleAuthFunc() {
	var _ AuthFunc = func(req *http.Request) bool {
		_, err := req.Cookie("authorized")
		return err == nil
	}
}

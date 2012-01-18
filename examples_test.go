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
	a.Init()

	http.Handle("/admin/", a)
	if err := http.ListenAndServe(":11223", nil); err != nil {
		log.Fatal(err)
	}
}

func ExampleAdmin() {
	session, err := mgo.Mongo("mongo://my_mongo_server")
	if err != nil {
		log.Fatal(err)
	}

	a := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
		Routes: map[string]string{
			"index":  "/1/",
			"list":   "/2/",
			"update": "/3/",
			"create": "/4/",
			"detail": "/5/",
			"delete": "/6/",
			"auth":   "/7/",
		},
		Prefix: "/admin",
	}

	a.Register(T{}, "database.collection", &Options{
		Columns: []string{"First", "Second", "Fifth"},
	})

	a.Init()

	http.Handle("/admin/", a)
	if err := http.ListenAndServe(":11223", nil); err != nil {
		log.Fatal(err)
	}
}

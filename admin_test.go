package admin

import (
	"net"
	"net/http"
	"testing"
)

type T struct{}

var conn chan net.Listener

func init() {
	conn = make(chan net.Listener)
	go func() {
		for {
			l, err := net.Listen("tcp", ":57222")
			if err != nil {
				panic("Unable to listen: " + err.Error())
			}

			conn <- l
			(<-conn).Close()
		}
	}()
}

func TestRegister(t *testing.T) {
	Register(T{}, "T", nil)

	ans := GetType("T")
	if _, ok := ans.(*T); !ok {
		t.Fatalf("Type incorrect. Expected *admin.T, got %T", ans)
	}
}

func TestUnauthorized(t *testing.T) {
	s := http.Server{
		Handler: &Admin{
			Auth: func(*http.Request) bool { return false },
		},
	}

	l := <-conn
	defer func() { conn <- l }()
	go s.Serve(l)

	resp, err := http.Get("http://localhost:57222/")
	if err != nil {
		t.Fatal("Failed getting response:", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatal("Failed being unauthorized")
	}
}

func TestAuthorized(t *testing.T) {
	s := http.Server{
		Handler: &Admin{
			Auth: func(*http.Request) bool { return true },
		},
	}

	l := <-conn
	defer func() { conn <- l }()
	go s.Serve(l)

	resp, err := http.Get("http://localhost:57222/")
	if err != nil {
		t.Fatal("Failed getting response:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Failed being unauthorized")
	}
}

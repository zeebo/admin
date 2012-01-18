package admin

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAuthRedirects(t *testing.T) {
	h := &Admin{
		Session: session,
		Auth:    TestAuth{},
	}
	h.Init()
	var w *TestResponseWriter

	w = Get(t, h, "/")
	if w.Status != http.StatusTemporaryRedirect {
		t.Fatalf("Expected %d. Got %d", http.StatusTemporaryRedirect, w.Status)
	}
}

func TestAuthLogin(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session: session,
		Auth: TestAuth{
			AuthResponse{
				Passed:   true,
				Username: "zeebo",
				Key:      "zeebos-id",
			},
		},
		Renderer: r,
	}
	h.Init()
	var w *TestResponseWriter

	w = Post(t, h, "/auth/login", url.Values{})
	cookie := w.Headers.Get("Set-Cookie")

	if !strings.HasPrefix(cookie, "auth=") {
		t.Fatal("Expected to get an auth cookie. Got %q", cookie)
	}
	if w.Status != http.StatusOK {
		t.Fatalf("[login] Expected status %d. Got %d.", http.StatusOK, w.Status)
	}

	//check the context stuff
	ctx, ok := r.Last().Params.(AuthorizeContext)
	if !ok {
		t.Fatalf("Expected AuthorizeContext. Got %T", r.Last().Params)
	}
	if ctx.Success != true {
		t.Fatalf("Success: Expected %v. Got %v", true, ctx.Success)
	}
	if ctx.Error != "" {
		t.Fatalf("Error: Expected %q. Got %q", "", ctx.Error)
	}

	chunks := strings.SplitN(cookie, "=", 2)
	data := strings.Split(chunks[1], ";")

	//have to make a request
	w = NewTestResponseWriter()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("unable to make request: %s", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: data[0],
	})

	h.ServeHTTP(w, req)
	w.Cleanup()

	if w.Status != http.StatusOK {
		t.Fatalf("[authd] Expected status %d. Got %d", http.StatusOK, w.Status)
	}
	if n := r.Last().Type; n != "Index" {
		t.Fatalf("Expected Index. Got %d", n)
	}
}

func TestAuthRedirectAfterLogin(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session: session,
		Auth: TestAuth{
			AuthResponse{
				Passed:   true,
				Username: "zeebo",
				Key:      "zeebos-id",
			},
		},
		Renderer: r,
	}
	h.Init()
	var w *TestResponseWriter

	w = NewTestResponseWriter()
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal("unable to make request: %s", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "redirect",
		Value: "/foo/bar",
	})

	h.ServeHTTP(w, req)
	w.Cleanup()

	if w.Status != http.StatusMovedPermanently {
		t.Fatalf("Expected %d. Got %d", http.StatusMovedPermanently, w.Status)
	}
	if redir := w.Headers.Get("Location"); redir != "/foo/bar" {
		t.Fatalf("Expected %q. Got %q", "/foo/bar", redir)
	}
}

func TestAuthRedirectHasWithPrefix(t *testing.T) {
	h := &Admin{
		Session: session,
		Auth:    TestAuth{},
		Prefix:  "/some/prefix",
	}
	h.Init()
	var w *TestResponseWriter

	w = Get(t, h, "/some/prefix/foo/bar")
	if w.Status != http.StatusTemporaryRedirect {
		t.Fatalf("Expected %d. Got %d", http.StatusTemporaryRedirect, w.Status)
	}
	cookie := w.Headers.Get("Set-Cookie")
	chunks := strings.SplitN(cookie, "=", 2)
	data := strings.Split(chunks[1], ";")

	if data[0] != "/some/prefix/foo/bar" {
		t.Fatalf("Expected %q. Got %q", "/some/prefix/foo/bar", data[0])
	}
}

func TestAuthLogout(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session: session,
		Auth: TestAuth{
			AuthResponse{
				Passed:   true,
				Username: "zeebo",
				Key:      "zeebos-id",
			},
		},
		Renderer: r,
	}
	h.Init()
	var w *TestResponseWriter

	w = Get(t, h, "/auth/logout")
	if cookie := w.Headers.Get("Set-Cookie"); !strings.HasPrefix(cookie, "auth=;") {
		t.Fatalf("Expected %q. Got %q", "auth=;*", cookie)
	}
}

func TestAuthFailedLoginErrors(t *testing.T) {
	r := &TestRenderer{}
	h := &Admin{
		Session: session,
		Auth: TestAuth{
			AuthResponse{
				Passed: false,
				Error:  "foo bar",
			},
		},
		Renderer: r,
	}
	h.Init()
	var w *TestResponseWriter

	w = Post(t, h, "/auth/login", url.Values{})
	if w.Status != http.StatusOK {
		t.Fatalf("Expected %d. Got %d", http.StatusOK, w.Status)
	}

	//make sure no auth cookie
	if c := w.Headers.Get("Set-Cookie"); c != "" {
		t.Fatalf("Expected %q. Got %q", "", c)
	}

	ctx, ok := r.Last().Params.(AuthorizeContext)
	if !ok {
		t.Fatalf("Expected AuthorizeContext. Got %T", r.Last().Params)
	}

	if ctx.Error != "foo bar" {
		t.Fatalf("Expected %q. Got %q.", "foo bar", ctx.Error)
	}
}

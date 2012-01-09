package admin

import (
	"net/http"
	"testing"
)

type doNothingResponseWriter struct{}

func (d doNothingResponseWriter) Header() http.Header {
	return make(http.Header)
}
func (d doNothingResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
func (d doNothingResponseWriter) WriteHeader(n int) {

}

func BenchmarkIndex(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	req, _ := http.NewRequest("GET", "/", nil)
	w := doNothingResponseWriter{}

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkList(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	req, _ := http.NewRequest("GET", "/list/admin_test.T/", nil)
	w := doNothingResponseWriter{}

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkUpdate(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	req, _ := http.NewRequest("GET", "/update/admin_test.T/4f07c34779bf562daff8640c", nil)
	w := doNothingResponseWriter{}

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkCreate(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	req, _ := http.NewRequest("GET", "/create/admin_test.T/", nil)
	w := doNothingResponseWriter{}

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkDetail(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	req, _ := http.NewRequest("GET", "/detail/admin_test.T/4f07c34779bf562daff8640c", nil)
	w := doNothingResponseWriter{}

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

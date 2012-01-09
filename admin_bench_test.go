package admin

import "testing"

func BenchmarkIndex(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	for i := 0; i < b.N; i++ {
		Request(h, "GET", "/", nil)
	}
}

func BenchmarkList(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	for i := 0; i < b.N; i++ {
		Request(h, "GET", "/list/admin_test.T/", nil)
	}
}

func BenchmarkUpdate(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	for i := 0; i < b.N; i++ {
		Request(h, "GET", "/update/admin_test.T/4f07c34779bf562daff8640c", nil)
	}
}

func BenchmarkCreate(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	for i := 0; i < b.N; i++ {
		Request(h, "GET", "/create/admin_test.T/", nil)
	}
}

func BenchmarkDetail(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.Register(T{}, "admin_test.T", nil)
	for i := 0; i < b.N; i++ {
		Request(h, "GET", "/detail/admin_test.T/4f07c34779bf562daff8640c", nil)
	}
}

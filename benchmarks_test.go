package admin

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
)

type doNothingResponseWriter struct{}

func (d doNothingResponseWriter) Header() http.Header         { return make(http.Header) }
func (d doNothingResponseWriter) Write(b []byte) (int, error) { return len(b), nil }
func (d doNothingResponseWriter) WriteHeader(n int)           {}

func channelSendHelper(p, q chan bool) {
	select {
	case q <- true:
	default:
		p <- true
		return
	}
	panic("never reached")
}

func BenchmarkChannelSend(b *testing.B) {
	var (
		p chan bool
		q = make(chan bool, 1)
	)
	q <- true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p = make(chan bool)
		go channelSendHelper(p, q)
		<-p
	}
}

func BenchmarkReverse(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	r := Reverser{h}
	h.Register(T{}, "admin_test.T", nil)
	h.generateMux()
	var x T = T{"ffffffffffff"}

	for i := 0; i < b.N; i++ {
		r.DetailObj(x)
	}
}

func BenchmarkGetIndex(b *testing.B) {
	h := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
	}
	h.Register(T{}, "admin_test.T", nil)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatal(err)
	}
	w := doNothingResponseWriter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkGetDelete(b *testing.B) {
	h := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
	}
	h.Register(T{}, "admin_test.T", nil)
	req, err := http.NewRequest("GET", "/delete/admin_test.T/4f07c34779bf562daff8640c", nil)
	if err != nil {
		b.Fatal(err)
	}
	w := doNothingResponseWriter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkGetList(b *testing.B) {
	h := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
	}
	h.Register(T{}, "admin_test.T", nil)
	req, err := http.NewRequest("GET", "/list/admin_test.T/", nil)
	if err != nil {
		b.Fatal(err)
	}
	w := doNothingResponseWriter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkGetUpdate(b *testing.B) {
	h := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
	}
	h.Register(T{}, "admin_test.T", nil)
	req, err := http.NewRequest("GET", "/update/admin_test.T/4f07c34779bf562daff8640c", nil)
	if err != nil {
		b.Fatal(err)
	}
	w := doNothingResponseWriter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkPostUpdate(b *testing.B) {
	h := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
	}
	h.Register(T6{}, "admin_test.T6", nil)
	w := doNothingResponseWriter{}
	var (
		values url.Values
		req    *http.Request
		err    error
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		values = url.Values{
			"X": {fmt.Sprint(rand.Intn(1000))},
			"Y": {"foo"},
			"Z": {"true"},
		}
		req, err = http.NewRequest("POST", "/update/admin_test.T6/4f0ee3600888a1b6646199bd", bytes.NewBufferString(values.Encode()))
		if err != nil {
			b.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		b.StartTimer()

		h.ServeHTTP(w, req)
	}
}

func BenchmarkGetCreate(b *testing.B) {
	h := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
	}
	h.Register(T{}, "admin_test.T", nil)
	req, err := http.NewRequest("GET", "/create/admin_test.T/", nil)
	if err != nil {
		b.Fatal(err)
	}
	w := doNothingResponseWriter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkGetDetail(b *testing.B) {
	h := &Admin{
		Session:  session,
		Renderer: &TestRenderer{},
	}
	h.Register(T{}, "admin_test.T", nil)
	req, err := http.NewRequest("GET", "/detail/admin_test.T/4f07c34779bf562daff8640c", nil)
	if err != nil {
		b.Fatal(err)
	}
	w := doNothingResponseWriter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkCRUDCycle(b *testing.B) {
	r := &TestRenderer{}
	h := &Admin{
		Session:  session,
		Renderer: r,
	}
	h.Register(T6{}, "admin_test.T6", nil)

	var id string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Post(b, h, "/create/admin_test.T6/", url.Values{
			"X": {"2"},
			"Y": {"new"},
			"Z": {"false"},
		})
		id = r.Last().Params.(CreateContext).Form.context.Values["ID"]

		Get(b, h, fmt.Sprintf("/detail/admin_test.T6/%s", id))
		Post(b, h, fmt.Sprintf("/update/admin_test.T6/%s", id), url.Values{
			"X": {"20"},
			"Y": {"newt"},
			"Z": {"true"},
		})
		Get(b, h, fmt.Sprintf("/delete/admin_test.T6/%s?_sure=yes", id))
	}
}

func BenchmarkAdminEmptyInit(b *testing.B) {
	h := &Admin{
		Session: session,
	}
	h.init()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.init()
	}
}

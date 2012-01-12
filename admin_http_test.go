package admin

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type TestResponseWriter struct {
	Body        bytes.Buffer
	Headers     http.Header
	Status      int
	WroteHeader bool
}

func (t *TestResponseWriter) Header() http.Header {
	return t.Headers
}

func (t *TestResponseWriter) Write(b []byte) (int, error) {
	return t.Body.Write(b)
}

func (t *TestResponseWriter) WriteHeader(n int) {
	if !t.WroteHeader {
		t.Status = n
		t.WroteHeader = true
	}
}
func (t *TestResponseWriter) Cleanup() {
	if !t.WroteHeader {
		t.WriteHeader(http.StatusOK)
	}
}

func NewTestResponseWriter() *TestResponseWriter {
	return &TestResponseWriter{
		Headers: make(http.Header),
	}
}

type fatalf interface {
	Fatalf(string, ...interface{})
}

func Request(h http.Handler, method string, url, bodyType string, body io.Reader) (*TestResponseWriter, error) {
	w := NewTestResponseWriter()
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if bodyType != "" {
		r.Header.Add("Content-Type", bodyType)
	}
	h.ServeHTTP(w, r)
	w.Cleanup()
	return w, nil
}

func Get(t fatalf, h http.Handler, url string) *TestResponseWriter {
	w, err := Request(h, "GET", url, "", nil)

	if err != nil {
		t.Fatalf("Error requesting %q: %s", url, err)
	}

	return w
}

func Post(t fatalf, h http.Handler, url string, data url.Values) *TestResponseWriter {
	buf := bytes.NewBufferString(data.Encode())
	w, err := Request(h, "POST", url, "application/x-www-form-urlencoded", buf)

	if err != nil {
		t.Fatalf("Error requesting %q: %s", url, err)
	}

	return w
}

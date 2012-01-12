package admin

import (
	"log"
	"net/http"
)

type TestCall struct {
	Type   string
	Params interface{}
}

type TestRenderer struct {
	Calls []TestCall
}

func (r *TestRenderer) Last() TestCall {
	if r.Calls == nil {
		return TestCall{"None", nil}
	}
	return r.Calls[len(r.Calls)-1]
}

func (r *TestRenderer) NotFound(w http.ResponseWriter, req *http.Request) {
	r.Calls = append(r.Calls, TestCall{
		Type: "NotFound",
	})
	w.WriteHeader(http.StatusNotFound)
}

func (r *TestRenderer) InternalError(w http.ResponseWriter, req *http.Request, err error) {
	r.Calls = append(r.Calls, TestCall{
		Type:   "InternalError",
		Params: err,
	})
	log.Println("Internal:", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (r *TestRenderer) Unauthorized(w http.ResponseWriter, req *http.Request) {
	r.Calls = append(r.Calls, TestCall{
		Type: "Unauthorized",
	})
}

func (r *TestRenderer) Detail(w http.ResponseWriter, req *http.Request, c DetailContext) {
	r.Calls = append(r.Calls, TestCall{
		Type:   "Detail",
		Params: c,
	})
}

func (r *TestRenderer) Delete(w http.ResponseWriter, req *http.Request, c DeleteContext) {
	r.Calls = append(r.Calls, TestCall{
		Type:   "Delete",
		Params: c,
	})
}

func (r *TestRenderer) Index(w http.ResponseWriter, req *http.Request, c IndexContext) {
	r.Calls = append(r.Calls, TestCall{
		Type:   "Index",
		Params: c,
	})
}

func (r *TestRenderer) List(w http.ResponseWriter, req *http.Request, c ListContext) {
	r.Calls = append(r.Calls, TestCall{
		Type:   "List",
		Params: c,
	})
}

func (r *TestRenderer) Update(w http.ResponseWriter, req *http.Request, c UpdateContext) {
	r.Calls = append(r.Calls, TestCall{
		Type:   "Update",
		Params: c,
	})
}

func (r *TestRenderer) Create(w http.ResponseWriter, req *http.Request, c CreateContext) {
	r.Calls = append(r.Calls, TestCall{
		Type:   "Create",
		Params: c,
	})
}

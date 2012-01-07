package admin

import "net/http"

type Renderer interface {
	//Error modes
	//
	//NotFound must return a http.StatusNotFound and InternalError must return
	//an http.StatusInternalServiceError to be compliant. The error that caused
	//the exception is passed in.
	NotFound(http.ResponseWriter, *http.Request)
	InternalError(http.ResponseWriter, *http.Request, error)
	Unauthorized(http.ResponseWriter, *http.Request)

	//Handler modes
	//
	//These functions present the user some html with the appropriate context
	//data.
	//
	//In the case of Update and Create, these handlers must handle the case
	//of an error in processing the form. Relevant details will be in the 
	//passed in context.
	Detail(http.ResponseWriter, *http.Request, DetailContext)
	Index(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request, ListContext)
	Update(http.ResponseWriter, *http.Request, UpdateContext)
	Create(http.ResponseWriter, *http.Request, CreateContext)
}

type DetailContext struct {
}

type ListContext struct {
}

type UpdateContext struct {
}

type CreateContext struct {
}

type DefaultRenderer struct{}

func (r DefaultRenderer) NotFound(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}

func (r DefaultRenderer) InternalError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (r DefaultRenderer) Unauthorized(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (r DefaultRenderer) Detail(w http.ResponseWriter, req *http.Request, c DetailContext) {

}

func (r DefaultRenderer) Index(w http.ResponseWriter, req *http.Request) {

}

func (r DefaultRenderer) List(w http.ResponseWriter, req *http.Request, c ListContext) {

}

func (r DefaultRenderer) Update(w http.ResponseWriter, req *http.Request, c UpdateContext) {

}

func (r DefaultRenderer) Create(w http.ResponseWriter, req *http.Request, c CreateContext) {

}

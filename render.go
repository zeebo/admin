package admin

import (
	"fmt"
	"net/http"
)

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
	Index(http.ResponseWriter, *http.Request, IndexContext)
	List(http.ResponseWriter, *http.Request, ListContext)
	Update(http.ResponseWriter, *http.Request, UpdateContext)
	Create(http.ResponseWriter, *http.Request, CreateContext)
}

//DetailContext is the type passed to the Detail method.
//It comes loaded with the instance of the object found. If the object cannot
//be found, the value wil be nil.
type DetailContext struct {
	Object interface{}
}

//ListContext is the type passed in to the List method.
//It comes loaded with a slice of objects selected by the List view. If no
//objects match the passed in query, the slice will be nil.
type ListContext struct {
	Objects []interface{}
}

//UpdateContext is the type passed in to the Update method.
//It comes with a map[string]string of Field -> Error values if there were any
//in the processing of the Update request. If there were no errors then 
//the map will be nil. It also comes with an instance of the object with the 
//matching query. If there is no matching object, it will be nil. The object
//always reflects the most recent data in the database. 
type UpdateContext struct {
	Object interface{}
	Errors map[string]string
}

//CreateContext is the type passed in to the Create method.
//It comes with a map[string]string of Field -> Error vlaues if there were any
//in the processing of the Create request. If there were no erros then the
//map will be nil.
type CreateContext struct {
	Errors map[string]string
}

//IndexContext is the type passed in to the Index method. It contains the
//databases and collections being managed by the admin.
type IndexContext struct {
	Managed map[string][]string
}

//Key takes a database and collection and maps it to the key for urls. For
//example, Key("db", "coll") -> db.coll
func (i *IndexContext) Key(db, coll string) string {
	return fmt.Sprintf("%s.%s", db, coll)
}

//DefaultRenderer conforms to the Renderer interface and uses some magic templates
//to create a pretty default interface.
type DefaultRenderer struct{}

//NotFound presents a basic 404 with no special body.
func (r DefaultRenderer) NotFound(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}

//InternalError presents a basic 500 not suitable for production. Errors should be logged
//and not displayed to the end user.
func (r DefaultRenderer) InternalError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

//Unauthorized presents a simple unauthorized page.
func (r DefaultRenderer) Unauthorized(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

//Detail presents the detail view of an object.
func (r DefaultRenderer) Detail(w http.ResponseWriter, req *http.Request, c DetailContext) {

}

//Index presents an overall view of the database and the managed collections.
func (r DefaultRenderer) Index(w http.ResponseWriter, req *http.Request, c IndexContext) {

}

//List presents all of the objects of a specific list with the columns and order given by the options
//the type was loaded with.
func (r DefaultRenderer) List(w http.ResponseWriter, req *http.Request, c ListContext) {

}

//Update presents a success page or the errors when attempting to update an object.
func (r DefaultRenderer) Update(w http.ResponseWriter, req *http.Request, c UpdateContext) {

}

//Create presents a success page or the errors when attempting to create an object.
func (r DefaultRenderer) Create(w http.ResponseWriter, req *http.Request, c CreateContext) {

}

package admin

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

//Renderer represents a type that knows how to present content for the admin to
//the browser.
type Renderer interface {
	//Error modes
	//
	//NotFound must return a http.StatusNotFound, InternalError must return
	//an http.StatusInternalServiceError. The error that caused the exception is
	//passed in.
	NotFound(http.ResponseWriter, *http.Request)
	InternalError(http.ResponseWriter, *http.Request, error)

	//Handler modes
	//
	//These functions present the user some html with the appropriate context
	//data.
	//
	//In the case of Update, Create and Authorize, these handlers must handle
	//the case of an error in processing the form. Relevant details will be in
	//the  passed in context.
	Detail(http.ResponseWriter, *http.Request, DetailContext)
	Delete(http.ResponseWriter, *http.Request, DeleteContext)
	Index(http.ResponseWriter, *http.Request, BaseContext)
	List(http.ResponseWriter, *http.Request, ListContext)
	Update(http.ResponseWriter, *http.Request, UpdateContext)
	Create(http.ResponseWriter, *http.Request, CreateContext)
	Authorize(http.ResponseWriter, *http.Request, AuthorizeContext)
	LoggedOut(http.ResponseWriter, *http.Request, BaseContext)
}

//DetailContext is the type passed to the Detail method.
//It comes loaded with the instance of the object found, and a Form that
//represents the form for the object.
type DetailContext struct {
	BaseContext
	Collection string
	Object     interface{}
	Form       Form
}

//DeleteContext is the type passed to the Delete method.
//It comes loaded with an instance of the object to be deleted, and a form
//for rendering the object. The renderer should use the Form.Values method to
//render a readonly display, with a button that adds _sure=yes as a parameter
//to the same page. Error is the error in attempting to delete the object, if
//one exists.
type DeleteContext struct {
	BaseContext
	Collection string
	Object     interface{}
	Attempted  bool
	Success    bool
	Error      error
	Form       Form
}

//ListContext is the type passed in to the List method.
//It comes loaded with a slice of objects selected by the List view. If no
//objects match the passed in query, the slice will be nil.
type ListContext struct {
	BaseContext
	Collection string
	Columns    []string
	Values     [][]string
	Objects    []interface{}
	Pagination Pagination
}

//UpdateContext is the type passed in to the Update method.
//It comes with booleans indicating if the update was attempted and successful.
//It also comes with an instance of the object with the matching query.
//The object always reflects the most recent data in the database.
//It also comes with a Form that represents the form for the object.
type UpdateContext struct {
	BaseContext
	Collection string
	Object     interface{}
	Attempted  bool
	Success    bool
	Error      error
	Form       Form
}

//CreateContext is the type passed in to the Create method.
//It comes with booleans indicating if the creation was attempted and successful.
//It also comes with a Form that represents the form for the object.
type CreateContext struct {
	BaseContext
	Collection string
	Attempted  bool
	Success    bool
	Error      error
	Form       Form
}

//AuthorizeContext is the type passed in to the Authorize method.
//It comes with booleans indicating if the authorization request was attempted
//and sucessful. It also comes with an error string if not sucessful.
type AuthorizeContext struct {
	BaseContext
	Success   bool
	Attempted bool
	Error     string
}

//BaseContext is the type passed in to every Render method. It contains the
//databases and collections being managed by the admin and information regarding
//the logged in user.
type BaseContext struct {
	Managed  map[string][]string
	Reverser Reverser
	Auth     *AuthSession
}

//Key takes a database and collection and maps it to the key for urls. For
//example, Key("db", "coll") -> db.coll
func (i BaseContext) Key(db, coll string) string {
	return fmt.Sprintf("%s.%s", db, coll)
}

//TemplateContext is the value passed in as the dot to the template for forms
//by the default renderer. It has methods for returning the values in the field
//and any errors in attempting to validate the form.
type TemplateContext struct {
	Errors map[string]interface{}
	Values map[string]interface{}
}

//NewTemplateContext creates a new TemplateContext ready to be used.
func NewTemplateContext() TemplateContext {
	return TemplateContext{map[string]interface{}{}, map[string]interface{}{}}
}

//Form encapsulates a form with a context with the ability to execute and output
//the correct html.
type Form struct {
	object  Formable
	context TemplateContext
	logger  *log.Logger
}

//Execute calls the template with the context and executes it to the writer
func (f Form) Execute(w io.Writer) (err error) {
	_, err = io.WriteString(w, f.object.GetForm(f.context))
	return
}

//ExecuteText is for use in templates. It returns the string containing the
//output of Execute.
func (f Form) ExecuteText() string {
	return f.object.GetForm(f.context)
}

//Values returns the values map for the Form. This is useful for the Delete
//renderer to display readonly values for the form.
func (f Form) Values() map[string]interface{} {
	return f.context.Values
}

package admin

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
)

//Renderer represents a type that knows how to present content for the admin to
//the browser.
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
	Delete(http.ResponseWriter, *http.Request, DeleteContext)
	Index(http.ResponseWriter, *http.Request, IndexContext)
	List(http.ResponseWriter, *http.Request, ListContext)
	Update(http.ResponseWriter, *http.Request, UpdateContext)
	Create(http.ResponseWriter, *http.Request, CreateContext)
}

//DetailContext is the type passed to the Detail method.
//It comes loaded with the instance of the object found, and a Form that
//represents the form for the object.
type DetailContext struct {
	Object interface{}
	Form   Form
}

//DeleteContext is the type passed to the Delete method.
//It comes loaded with an instance of the object to be deleted, and a form
//for rendering the object. The renderer should use the Form.Values method to
//render a readonly display, with a button that adds _sure=yes as a parameter
//to the same page. Error is the error in attempting to delete the object, if
//one exists.
type DeleteContext struct {
	Object    interface{}
	Attempted bool
	Success   bool
	Error     error
	Form      Form
}

//ListContext is the type passed in to the List method.
//It comes loaded with a slice of objects selected by the List view. If no
//objects match the passed in query, the slice will be nil.
type ListContext struct {
	Objects []interface{}
}

//UpdateContext is the type passed in to the Update method.
//It comes with booleans indicating if the update was attempted and successful.
//It also comes with an instance of the object with the matching query.
//The object always reflects the most recent data in the database.
//It also comes with a Form that represents the form for the object.
type UpdateContext struct {
	Object    interface{}
	Attempted bool
	Success   bool
	Form      Form
}

//CreateContext is the type passed in to the Create method.
//It comes with booleans indicating if the creation was attempted and successful.
//It also comes with a Form that represents the form for the object.
type CreateContext struct {
	Attempted bool
	Success   bool
	Form      Form
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

//TemplateContext is the value passed in as the dot to the template for forms
//by the default renderer. It has methods for returning the values in the field
//and any errors in attempting to validate the form. For example if we had the
//type
//
//	type MyForm struct {
//		X int
//		Y string
//	}
//
//a simple template that uses the TemplateContext for this struct could look like
//
//	func (m *MyForm) GetTemplate() string {
//		return `<span class="errors">{{.Errors "X"}}</span>
//			<input type="text" value="{{.Values "X"}}" name="X">
//			<span class="errors">{{.Errors "Y"}}</span>
//			<input type="text" value="{{.Values "Y"}}" name="Y">
//			<input type="submit">`
//	}
//
//The form is rendered through the html/template package and will do necessary
//escaping as such. It is the renderers responsibility to wrap the fields
//in a form tag.
type TemplateContext struct {
	Errors map[string]error
	Values map[string]string
}

//NewTemplateContext creates a new TemplateContext ready to be used.
func NewTemplateContext() TemplateContext {
	return TemplateContext{map[string]error{}, map[string]string{}}
}

//Error returns any error text from validation for a specific field.
func (t *TemplateContext) Error(field string) string {
	if v, ex := t.Errors[field]; ex && v != nil {
		return v.Error()
	}
	return ""
}

//Value returns the value the user input into the form after validation.
func (t *TemplateContext) Value(field string) string {
	if v, ex := t.Values[field]; ex {
		return v
	}
	return ""
}

//Form encapsulates a form with a context with the ability to execute and output
//the correct html.
type Form struct {
	template *template.Template
	context  TemplateContext
}

//Execute calls the template with the context and executes it to the writer
func (f *Form) Execute(w io.Writer) error {
	return f.template.Execute(w, f.context)
}

//Values returns the values map for the Form. This is useful for the Delete
//renderer to display readonly values for the form.
func (f *Form) Values() map[string]string {
	return f.context.Values
}

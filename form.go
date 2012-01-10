package admin

//Formable is the type of objects that the admin can represent.
type Formable interface {
	//GetTemplate returns the template text that will be used to render the
	//form on the page. For more details on what gets sent into the form
	//and what methods are present for rendering, see TemplateContext
	GetTemplate() string

	//Validate is called on the type after all the individual fields have been
	//validated and no errors have occured.
	Validate() error
}

//TemplateContext is the value passed in as the dot to the template for forms.
//It has methods for returning the values in the field and any errors in
//attempting to validate the form. For example if we had the struct
//
//	type MyForm struct {
//		X int
//		Y string
//	}
//
//a simple template that uses the TemplateContext for this struct could look like
//
//	func (m *MyForm) GetTemplate() string {
//		return `
//			<form method="post" action="{{.Action}}">
//				<span class="errors">{{.Errors "X"}}</span>
//				<input type="text" value="{{.Values "X"}}" name="X">
//				<span class="errors">{{.Errors "Y"}}</span>
//				<input type="text" value="{{.Values "Y"}}" name="Y">
//				<input type="submit">
//			</form>`
//	}
//
//The form is rendered through the html/template package and will do necessary
//escaping as such.
type TemplateContext struct {
	Errors map[string]string
	Values map[string]string
	Action string
}

//Error returns any error text from validation for a specific field.
func (t *TemplateContext) Error(field string) string {
	if v, ex := t.Errors[field]; ex {
		return v
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

//Validator represents a type that can be validated by the form processor. For
//example, we can make a string field that cannot have numbers like
//
//	type NoNumberField string
//	func (n *NoNumberField) Validate() error {
//		if strings.IndexAny(*n, "0123456789") != -1 {
//			return errors.New("This field must contain no numbers.")
//		}
//		return nil
//	}
//
//The form processor will check if the type of the field is a Validator and do
//any validation required. This method is allowed to change the data for the
//field.
type Validator interface {
	Validate() error
}

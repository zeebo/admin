package admin

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

//Formable is the type of objects that the admin can represent.
type Formable interface {
	//GetTemplate returns the template text that will be used to render the
	//form on the page. For more details on what gets sent into the form
	//and what methods are present for rendering, see TemplateContext
	GetTemplate() string

	//Validate is called on the type after all the individual fields are loaded.
	//There must be no errors loading for Validate to be called.
	Validate() ValidationErrors
}

//Loader allows you to define your own custom form loading methods if the built
//in automatic loading does not suit you. This method will be called instead of
//any other loading method if the type conforms to the interface. See 
//LoadingErrors for a description of how to generate this value. Note that if
//an inner struct has errors, you do not need to worry about the prefix. For
//example, it is ok, given the type
//
//	type T struct {
//		A string
//		B struct {
//			C string
//			D uint
//		}
//	}
//
//for B, if it is a Loader, to return LoadingErrors with the keys "C" and "D".
//The second paramater is used to indicate errors that don't have to do with
//loading such as an incorrect schema sent to your struct.
//
//GenerateContext returns the TemplateContext that will be passed in to the
//template returned by GetTemplate(). This cannot be handled by the admin
//if you are using types such as slices or maps in your data structure.
type Loader interface {
	Formable
	Load(url.Values) (LoadingErrors, error)
	GenerateContext() TemplateContext
}

//LoadingErrors is the type that the Load method returns for errors loading into
//a struct. For example trying to put "-1" into a uint or other things. Keys
//must be a dot seperated path to the value in the struct. For example
//
//	type T struct {
//		A string
//		B struct {
//			C string
//			D uint
//		}
//	}
//
//has the keys "A", "B.C", and "B.D".
type LoadingErrors map[string]error

//ValidationErrors is the type that Validate must return to indicate errors in
//validation. For example putting an invalid email or phone number into a field
//expecting one. For what the keys of the map must be, see LoadingErrors
type ValidationErrors map[string]error

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
//in the form.
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

//indirect walks up interface/pointer chains until it gets to an actual concrete
//type. If the pointer is nil, we can't walk up so we get an error.
func indirect(val reflect.Value) (v reflect.Value, e error) {
	//recover any errors from reflect
	defer func() {
		if i := recover(); i != nil {
			if err, ok := i.(error); ok {
				e = err
			}
			panic(i)
		}
	}()

	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if !val.IsValid() {
		return val, errors.New("Invalid value after indirection")
	}

	return val, nil
}

func indirectType(val reflect.Type) reflect.Type {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

//CreateValues is used to create a map for insertion into a TemplateContext.
//It calls fmt.Sprintf on the values which hopefully wont mangle anything.
func CreateValues(obj interface{}) (map[string]string, error) {
	val, err := indirect(reflect.ValueOf(obj))
	if err != nil {
		return nil, err
	}
	typ := val.Type()

	res := map[string]string{}
	for i := 0; i < val.NumField(); i++ {
		field, err := indirect(val.Field(i))
		if err != nil {
			return nil, err
		}
		name := typ.Field(i).Name

		if !field.CanInterface() {
			return nil, errors.New(fmt.Sprintf("Can't get the value in %s", name))
		}

		//handle the basic types
		if field.Kind() != reflect.Struct {
			res[name] = fmt.Sprint(field.Interface())
			continue
		}

		//handle the struct type
		data, err := CreateValues(field.Interface())
		if err != nil {
			return nil, err
		}

		//copy data into local map
		for k, v := range data {
			res[fmt.Sprintf("%s.%s", name, k)] = v
		}
	}

	return res, nil
}

//CreateEmptyValues creates a map for insertion into a TemplateContext using
//the empty string for every value. This is useful for generating a template
//context on a type that has not been loaded into, e.g. the create page.
func CreateEmptyValues(obj interface{}) (map[string]string, error) {
	typ := indirectType(reflect.TypeOf(obj))
	return createEmptyValuesType(typ)
}

//createEmptyValuesType is a helper for createEmptyValues. It helps the client
//not depend on the reflect package.
func createEmptyValuesType(typ reflect.Type) (m map[string]string, e error) {
	//capture errors because we're going cowboy with reflect
	defer func() {
		if i := recover(); i != nil {
			if err, ok := i.(error); ok {
				e = err
			}
			panic(i)
		}
	}()

	res := map[string]string{}
	for i := 0; i < typ.NumField(); i++ {
		field, name := indirectType(typ.Field(i).Type), typ.Field(i).Name

		if !validType(field) {
			return nil, errors.New(fmt.Sprintf("Unsupported type: %s", field.Kind()))
		}

		if field.Kind() != reflect.Struct {
			res[name] = ""
			continue
		}

		data, err := createEmptyValuesType(field)
		if err != nil {
			return nil, err
		}

		//copy the data in
		for k, v := range data {
			res[fmt.Sprintf("%s.%s", name, k)] = v
		}
	}

	return res, nil
}

//validType checks to see if the reflect.Type's Kind is a supported type. These types
//are the basic go types (int/string/etc.). More may be supported in the future.
func validType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map, reflect.Uintptr,
		reflect.Complex128, reflect.Complex64, reflect.Func, reflect.UnsafePointer,
		reflect.Ptr, reflect.Interface:
		return false
	}
	return true
}

//alloc walks up a type through indirections and interfaces allocating as needed
//until it gets to a concrete base type.
func alloc(v reflect.Value) reflect.Value {
	//TODO: not sure if this is needed. figure that out.
	/*
		if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
			v = v.Addr()
		}
	*/
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			v = v.Elem()
			continue
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

//loadInto takes a value, allocates and sets the value to the data inspecting
//the underlying type to do correct conversions.
func loadInto(val reflect.Value, data string) (e error) {
	//alloc all dem pointers
	val = alloc(val)

	if !val.IsValid() || !val.CanSet() {
		return errors.New("Value cannot be assigned to.")
	}

	//catch any panics from reflect and just return it as an error
	defer func() {
		if i := recover(); i != nil {
			if err, ok := i.(error); ok {
				e = err
			}
			panic(i)
		}
	}()

	switch val.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(data)
		if err != nil {
			return err
		}
		val.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return err
		}
		val.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return err
		}
		val.SetFloat(v)
	case reflect.String:
		val.SetString(data)
	default:
		return errors.New(fmt.Sprintf("Can't insert into a %s", val.Kind()))
	}

	return nil
}

//unflatten takes some url.Values and turns it into a nested datastructure. Each
//value can only be a d or a string. In other words
//
//	ret := unflatten(form, "")
//	switch ret["key"].(type) {
//	case string:
//		//do something with string
//	case d:
//		//do something with the dict
//	default:
//		panic("never happens")
//	}
//
//It unflattens based on the "." seperator and ignores every value past the
//first in the values. For example, for the url.Values given by
//
//	url.Values{"A": {"a"}, "B.C": {"c"}, "B.D": {"d"}}
//
//we will get the output dictionary given by
//
//	d{"A": "a", "B": d{"C": "c", "D": "d"}}
//
//Undefined behavior for degenerate keys such as ".." or ""
func unflatten(val url.Values, prefix string) d {
	ret := d{}
	for key, v := range val {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		if v == nil || len(v) == 0 {
			continue
		}
		item := v[0]
		key = key[len(prefix):]

		if _, ex := ret[key]; ex {
			continue
		}

		if strings.Contains(key, ".") {
			first := strings.Split(key, ".")[0]
			ret[first] = unflatten(val, fmt.Sprintf("%s%s.", prefix, first))
		} else {
			ret[key] = item
		}
	}
	return ret
}

//Load writes the values from a form into the specified object. Types handled
//are some basic types, types that are identical to those basic types, and 
//structs consisting of the former. If the object is a Loader loading is passed
//off to its Load method. If the types 
func Load(form url.Values, obj interface{}) (LoadingErrors, error) {
	if l, ok := obj.(Loader); ok {
		return l.Load(form)
	}

	val := reflect.ValueOf(obj).Elem()
	if !val.CanSet() || !val.IsValid() {
		return nil, errors.New(fmt.Sprintf("Can't set to the object sent in: CanSet(%v) IsValid(%v)", val.CanSet(), val.IsValid()))
	}

	return apply(val, unflatten(form, ""), "")
}

//apply does the heavy lifting for Load, recursing down the type when needed
//and mangling the LoadingErrors returned to have the correct prefix. obj.Kind()
//should always be reflect.Struct. Note prefix is only used to generate better
//error messages in the case of problems.
func apply(obj reflect.Value, data d, prefix string) (LoadingErrors, error) {
	//make sure we have a good value
	if obj.Kind() != reflect.Struct || !obj.CanSet() || !obj.IsValid() {
		return nil, errors.New(fmt.Sprintf("Attempted to apply on something that wasn't a struct or was invalid - CanSet(%v) IsValid(%v) Kind(%s)", obj.CanSet(), obj.IsValid(), obj.Kind()))
	}

	//set up our holders
	typ, errs := obj.Type(), LoadingErrors{}

	for i := 0; i < obj.NumField(); i++ {
		field, name := alloc(obj.Field(i)), typ.Field(i).Name
		if _, ex := data[name]; !ex {
			continue
		}

		//make sure the field is ok
		if t := indirectType(typ.Field(i).Type); !validType(t) {
			return nil, errors.New(fmt.Sprintf("Attempted to load into a %v, an invalid type.", t))
		}

		//handle basic field types
		if field.Kind() != reflect.Struct {
			sval, ok := data[name].(string)
			if !ok {
				return nil, errors.New("Attmped to load a dictionary into a basic type: " + prefix + name)
			}

			//load the thing into the field and grab the errors
			if err := loadInto(field, sval); err != nil {
				errs[name] = err
			}

			continue
		}

		//handle the struct case
		dval, ok := data[name].(d)
		if !ok {
			return nil, errors.New("Attempted to load a string into a struct type: " + prefix + name)
		}

		//recurse
		nest_err, ferr := apply(field, dval, fmt.Sprintf("%s%s.", prefix, name))
		if ferr != nil {
			return nil, ferr
		}

		//copy the nested errors into our map
		for key, err := range nest_err {
			errs[fmt.Sprintf("%s.%s", name, key)] = err
		}
	}

	return errs, nil
}

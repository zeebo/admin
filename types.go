package admin

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

//Options when adding a collection to the admin
type Options struct {
	//Which columns to display/order to display them - nil means all
	Columns []string
}

//findIds finds the index locations of the type matching the columns passed in.
func findIds(typ reflect.Type, columns []string) []int {
	//if columns is nil, do every column!
	if columns == nil {
		ids := make([]int, typ.NumField())
		for i := range ids {
			ids[i] = i
		}
		return ids
	}

	//otherwise
	ids := make([]int, len(columns))
	for i, col := range columns {
		var (
			found bool
			field reflect.StructField
		)

		//do a simple search
		for j := 0; j < typ.NumField(); j++ {
			field = typ.Field(j)
			if field.Name == col {
				ids[i] = j
				found = true
				break
			}
		}

		//panic if we didnt find it
		if !found {
			panic(fmt.Sprintf("Can't find a column named %s on type %s", col, typ))
		}
	}
	return ids
}

//Stores info about a specific collection, like the type of the object it
//represents and any options used in specifying the type
type collectionInfo struct {
	Type      reflect.Type
	Template  *template.Template
	ColumnIds []int
}

//Registers the type/collection pair in the admin. Panics if two types are mapped
//to the same collection. Dbcolls are dot separated database/collection specifiers.
//Panics if no database is specified. Panics if the template returned by the Formable
//has any compilation errors. Panics if the type cannot be handled by the loading
//engine (must be composed of valid types. See Load for discussion on which types
//are valid.) Panics if it can't find a field with a bson:_id tag.
func (a *Admin) Register(typ Formable, dbcoll string, opt *Options) {
	if a.types == nil {
		a.types = make(map[string]collectionInfo)
		a.object_id = make(map[reflect.Type]int)
		a.object_coll = make(map[reflect.Type]string)
	}

	if !strings.Contains(dbcoll, ".") {
		panic("Database/collection specifier does not contain a .")
	}
	t := indirectType(reflect.TypeOf(typ))
	if ci, ok := a.types[dbcoll]; ok {
		panic(fmt.Sprintf("db.collection already registered. Had %q->%s . Got %q->%s", dbcoll, ci.Type, dbcoll, t))
	}
	//compile the template
	templ := template.Must(template.New("form").Parse(typ.GetTemplate()))

	//ensure we can create an empty templatecontext for it (no invalid types)
	if _, err := CreateEmptyValues(typ); err != nil {
		//we have a type that cant be handled by the loading enging. But is it a
		//loader?
		if _, ok := typ.(Loader); !ok {
			panic(err)
		}
	}

	//now ensure that we can find out where the id is. Look for a bson:_id tag
	var i int
	for i = 0; i < t.NumField(); i++ {
		field := t.Field(i)
		for _, tag := range strings.Split(field.Tag.Get("bson"), ",") {
			if tag == "_id" {
				goto found
			}
		}
	}
	panic("Unable to find a field that is an id. Be sure to add a bson:_id to your struct")

found:
	//time to load up our data
	a.object_id[t] = i
	a.object_coll[t] = dbcoll

	var ids []int

	if opt == nil {
		ids = findIds(t, nil)
	} else {
		ids = findIds(t, opt.Columns)
	}

	a.types[dbcoll] = collectionInfo{t, templ, ids}
}

//Unregisters the information for the colleciton. Panics if you attempt to unregister
//a collection not yet registered.
func (a *Admin) Unregister(dbcoll string) {
	if a.types == nil {
		a.types = make(map[string]collectionInfo)
	}

	if _, ok := a.types[dbcoll]; !ok {
		panic(fmt.Sprintf("unregister db.collection that does not exist: %q", dbcoll))
	}
	delete(a.types, dbcoll)
}

//hasType returns if the database/collection pair has been registered.
func (a *Admin) hasType(dbcoll string) (ok bool) {
	if a.types == nil {
		return false
	}

	_, ok = a.types[dbcoll]
	return
}

//Returns an interface{} boxing a new(T) where T is the type registered
//under the collection name.
func (a *Admin) newType(dbcoll string) Formable {
	if a.types == nil {
		return nil
	}

	t, ok := a.types[dbcoll]
	if !ok {
		return nil
	}

	return reflect.New(t.Type).Interface().(Formable)
}

package admin

import (
	"fmt"
	"reflect"
	"strings"
)

//Options when adding a collection to the admin
type Options struct {
	//Which columns to display/order to display them - nil means all
	Columns []string
}

//Stores info about a specific collection, like the type of the object it
//represents and any options used in specifying the type
type collectionInfo struct {
	Options *Options
	Type    reflect.Type
}

//Registers the type/collection pair in the admin. Panics if two types are mapped
//to the same collection. Dbcolls are dot separated database/collection specifiers.
//Panics if no database is specified.
func (a *Admin) Register(typ interface{}, dbcoll string, opt *Options) {
	if a.types == nil {
		a.types = make(map[string]collectionInfo)
	}
	if !strings.Contains(dbcoll, ".") {
		panic("Database/collection specifier does not contain a .")
	}
	t := reflect.TypeOf(typ)
	if ci, ok := a.types[dbcoll]; ok {
		panic(fmt.Sprintf("db.collection already registered. Had %q->%s . Got %q->%s", dbcoll, ci.Type, dbcoll, t))
	}
	a.types[dbcoll] = collectionInfo{opt, t}
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
func (a *Admin) newType(dbcoll string) interface{} {
	if a.types == nil {
		return nil
	}

	t, ok := a.types[dbcoll]
	if !ok {
		return nil
	}

	return reflect.New(t.Type).Interface()
}

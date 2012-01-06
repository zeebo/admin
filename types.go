package admin

import (
	"fmt"
	"reflect"
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
//to the same collection
func (a *Admin) Register(typ interface{}, coll string, opt *Options) {
	if a.collections == nil {
		a.collections = make(map[string]collectionInfo)
	}

	t := reflect.TypeOf(typ)
	if ci, ok := a.collections[coll]; ok {
		panic(fmt.Sprintf("collection already registered. Had %q->%s . Got %q->%s", coll, ci.Type, coll, t))
	}
	a.collections[coll] = collectionInfo{opt, t}
}

//Unregisters the information for the colleciton. Panics if you attempt to unregister
//a collection not yet registered.
func (a *Admin) Unregister(coll string) {
	if a.collections == nil {
		a.collections = make(map[string]collectionInfo)
	}

	if _, ok := a.collections[coll]; !ok {
		panic(fmt.Sprintf("unregister collection that does not exist: %q", coll))
	}
	delete(a.collections, coll)
}

//Returns an interface{} boxing a new(T) where T is the type registered
//under the collection name.
func (a *Admin) newType(coll string) interface{} {
	if a.collections == nil {
		return nil
	}

	t, ok := a.collections[coll]
	if !ok {
		return nil
	}

	return reflect.New(t.Type).Interface()
}

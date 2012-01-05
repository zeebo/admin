package admin

import (
	"fmt"
	"reflect"
)

var collections = make(map[string]reflect.Type)

//Registers the type/collection pair in the admin. Panics if two types are mapped
//to the same collection
func Register(typ interface{}, collection string) {
	t := reflect.TypeOf(typ)
	if c, ok := collections[collection]; ok {
		panic(fmt.Sprintf("collection already registered: %s -> %s", c, t))
	}
	collections[collection] = t
}

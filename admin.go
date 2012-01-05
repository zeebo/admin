package admin

import (
	"fmt"
	"reflect"
)

var collections = make(map[string]reflect.Type)

func Register(type interface{}, collection string) error {
	typ := reflect.TypeOf(type)
	if c, ok = collections[collection]; ok {
		panic(fmt.Sprintf("collection already registered: %s -> %s", c, typ)
	}
	collections[collection] = typ
}
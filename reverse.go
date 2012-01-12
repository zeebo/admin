package admin

import (
	"fmt"
	"path"
	"reflect"
)

//Revserser is a type that allows turning data into urls.
type Reverser struct {
	admin *Admin
}

//idFor consults the admins object_id cache to find the field containing the ID
//for the object and returns that id as a string.
func (r Reverser) idFor(thing interface{}) string {
	typ, val := reflect.TypeOf(thing), reflect.ValueOf(thing)
	idx, ex := r.admin.object_id[typ]
	if !ex {
		panic(fmt.Sprintf("Don't know how to get the id for a %T", thing))
	}

	id := val.Field(idx).Interface()
	switch t := id.(type) {
	case hexable:
		return t.Hex()
	}
	return fmt.Sprint(id)
}

//collFor consults the admins object_coll cache to find the collection for the
//specified object's type.
func (r Reverser) collFor(thing interface{}) string {
	typ := indirectType(reflect.TypeOf(thing))
	col, ex := r.admin.object_coll[typ]
	if !ex {
		panic(fmt.Sprintf("Don't know how to get the collection for a %T", thing))
	}
	return col
}

func (r Reverser) CreateObj(thing interface{}) string {
	return r.Create(r.collFor(thing))
}
func (r Reverser) Create(coll string) string {
	return path.Join(r.admin.Prefix, r.admin.Routes["create"], coll)
}

func (r Reverser) DetailObj(thing interface{}) string {
	coll, id := r.collFor(thing), r.idFor(thing)
	return r.Detail(coll, id)
}
func (r Reverser) Detail(coll string, id string) string {
	return path.Join(r.admin.Prefix, r.admin.Routes["detail"], coll, id)
}

func (r Reverser) ListObj(thing interface{}) string {
	return r.List(r.collFor(thing))
}
func (r Reverser) List(coll string) string {
	return path.Join(r.admin.Prefix, r.admin.Routes["list"], coll)
}

func (r Reverser) Index() string {
	return path.Join(r.admin.Prefix, r.admin.Routes["index"])
}

func (r Reverser) DeleteObj(thing interface{}) string {
	coll, id := r.collFor(thing), r.idFor(thing)
	return r.Delete(coll, id)
}
func (r Reverser) Delete(coll string, id string) string {
	return path.Join(r.admin.Prefix, r.admin.Routes["delete"], coll, id)
}

func (r Reverser) UpdateObj(thing interface{}) string {
	coll, id := r.collFor(thing), r.idFor(thing)
	return r.Update(coll, id)
}
func (r Reverser) Update(coll string, id string) string {
	return path.Join(r.admin.Prefix, r.admin.Routes["update"], coll, id)
}

package admin

import (
	"fmt"
	"path"
	"reflect"
)

//Revserser is a type that allows turning data into urls. One is constructed for
//every context.
type Reverser struct {
	admin *Admin
}

//idFor consults the admins object_id cache to find the field containing the ID
//for the object and returns that id as a string.
func (r Reverser) idFor(thing interface{}) string {
	typ := indirectType(reflect.TypeOf(thing))
	val, err := indirect(reflect.ValueOf(thing))
	if err != nil {
		return ""
	}

	idx, ex := r.admin.object_id[typ]
	if !ex {
		panic(fmt.Errorf("Don't know how to get the id for a %T", thing))
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
		panic(fmt.Errorf("Don't know how to get the collection for a %T", thing))
	}
	return col
}

//CreateObj returns the url to create an object of the same type as the passed
//in object.
func (r Reverser) CreateObj(thing interface{}) string {
	r.admin.init()
	return r.Create(r.collFor(thing))
}

//Create returns the url to create an object of the given database/collection.
func (r Reverser) Create(coll string) string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["create"], coll)
}

//DetailObj returns the url to view the info of the passed in object.
func (r Reverser) DetailObj(thing interface{}) string {
	r.admin.init()
	coll, id := r.collFor(thing), r.idFor(thing)
	return r.Detail(coll, id)
}

//Detail returns the url to view the info of the passed in collection + object id.
func (r Reverser) Detail(coll string, id string) string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["detail"], coll, id)
}

//ListObj returns the url to view a list of objects with the same type as the 
//passed in object.
func (r Reverser) ListObj(thing interface{}) string {
	r.admin.init()
	return r.List(r.collFor(thing))
}

//List returns the url to view a list of objects for the given database/collection.
func (r Reverser) List(coll string) string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["list"], coll)
}

//Index returns the url for the admin index.
func (r Reverser) Index() string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["index"])
}

//DeleteObj returns the url to delete the passed in object.
func (r Reverser) DeleteObj(thing interface{}) string {
	r.admin.init()
	coll, id := r.collFor(thing), r.idFor(thing)
	return r.Delete(coll, id)
}

//Delete returns the url to delete the object given by the database/collection and id.
func (r Reverser) Delete(coll string, id string) string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["delete"], coll, id)
}

//UpdateObj returns the url to update the passed in object.
func (r Reverser) UpdateObj(thing interface{}) string {
	r.admin.init()
	coll, id := r.collFor(thing), r.idFor(thing)
	return r.Update(coll, id)
}

//Update returns the url to update the object given by database/collection and id.
func (r Reverser) Update(coll string, id string) string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["update"], coll, id)
}

func (r Reverser) Login() string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["auth"], "login")
}

func (r Reverser) Logout() string {
	r.admin.init()
	return path.Join(r.admin.Prefix, r.admin.Routes["auth"], "logout")
}

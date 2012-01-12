package admin

import (
	"launchpad.net/gobson/bson"
	"net/url"
)

//T is the most basic type possible
type T struct {
	ID bson.ObjectId `bson:"_id"`
}

func (t T) GetTemplate() string {
	return ``
}
func (t T) Validate() ValidationErrors { return nil }

//T2 is a type with data
type T2 struct {
	ID bson.ObjectId `bson:"_id"`
	V  int           `bson:"v"`
}

func (t T2) GetTemplate() string {
	return ``
}
func (t T2) Validate() ValidationErrors { return nil }

type T3 struct {
}

//T3 is a type that has an invalid template
func (t T3) GetTemplate() string {
	return `{{`
}
func (t T3) Validate() ValidationErrors { return nil }

//T4 is a type that cannot be managed by the loader
type T4 struct {
	x []string
}

func (t T4) GetTemplate() string {
	return ``
}
func (t T4) Validate() ValidationErrors { return nil }

//T5 is a type that cannot be managed by the loader but has a custom loader
//to allow it to work
type T5 struct {
	x []string
}

func (t T5) GetTemplate() string {
	return ``
}
func (t T5) Validate() ValidationErrors { return nil }

func (t T5) Load(v url.Values) (LoadingErrors, error) { return nil, nil }
func (t T5) GenerateContext() TemplateContext         { return NewTemplateContext() }

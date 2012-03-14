package admin

import (
	"launchpad.net/mgo/bson"
	"net/url"
)

//T is the most basic type possible
type T struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
}

func (t T) GetForm(ctx TemplateContext) string { return `` }
func (t T) Validate() ValidationErrors         { return nil }

var _ Formable = T{}

//T2 is a type with data
type T2 struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	V  int           `bson:"v"`
}

func (t T2) GetForm(ctx TemplateContext) string { return `` }
func (t T2) Validate() ValidationErrors         { return nil }

var _ Formable = T2{}

//T3 is a type that has an invalid template
type T3 struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
}

func (t T3) GetForm(ctx TemplateContext) string { return `` }
func (t T3) Validate() ValidationErrors         { return nil }

var _ Formable = T3{}

//T4 is a type that cannot be managed by the loader
type T4 struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	x  []string
}

func (t T4) GetForm(ctx TemplateContext) string { return `` }
func (t T4) Validate() ValidationErrors         { return nil }

var _ Formable = T4{}

//T5 is a type that cannot be managed by the loader but has a custom loader
//to allow it to work
type T5 struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	x  []string
}

func (t T5) GetForm(ctx TemplateContext) string       { return `` }
func (t T5) Validate() ValidationErrors               { return nil }
func (t T5) Load(v url.Values) (LoadingErrors, error) { panic("called l"); return nil, nil }
func (t T5) GenerateValues() map[string]interface{}   { panic("called gv"); return nil }

var _ Loader = T5{}
var _ Formable = T5{}

//T6 is a nontrivial type for testing CRUD
type T6 struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	X  int
	Y  string
	Z  bool
}

func (t T6) GetForm(ctx TemplateContext) string { return `` }
func (t T6) Validate() ValidationErrors         { return nil }

var _ Formable = T6{}

//T7 is a type that has no ID
type T7 struct{}

func (t T7) GetForm(ctx TemplateContext) string { return `` }
func (t T7) Validate() ValidationErrors         { return nil }

var _ Formable = T7{}

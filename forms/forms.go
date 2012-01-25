package forms

//Generator is a type that can generate the html representation of a field.
type Generator interface {
	Generate(Field, FieldContext) (string, error)
}

//A type is required to implement chooser when a field presents multiple choices.
type Chooser interface {
	Choices(string) []Item
}

//Item represents an item 
type Item struct {
	Label, Value string
}

//FieldContext is the type passed into the Generate method to determine how to
//generate the field. blah blah blah.
//TODO: come back when brain not stupid.
type FieldContext struct {
	Name  string
	Value interface{}
	Label string
	Error error

	//used for things like Radio/Select
	Choices []Item
}

//Only supports single valued fields at the moment
type Field string

//Enumerate the Fields we know how to thandle
const (
	Password Field = "Password"
	Text     Field = "Text"
	Textarea Field = "Textarea"
	Checkbox Field = "Checkbox"
	Radio    Field = "Radio"
	Select   Field = "Select"
)

func (f Field) String() string {
	return string(f)
}

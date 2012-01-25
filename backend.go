package admin

type Backend interface {
	Load(typ string, id interface{}, val interface{}) error
	Set(typ string, val interface{}) error
	List(typ string, spec ListSpec) ([]interface{}, error)
}

type ListSpec struct {
	NumPage int
	Page    int
	Sort    []SortType
}

type SortType struct {
	Field     string
	Direction SortDirection
}

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

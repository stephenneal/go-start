package ds

import (
	"sync"
)

var (
	write sync.Mutex
)

func NewModel0(code, name, colour, size string) Model0 {
	m := Model0{
		Name:   name,
		Code:   code,
	}
	return m
}

func ExistsModel0(code string) (bool, error) {
	c := col()
	return ExistsDoc(c, eqCode(code))
}

func FindModel0(code string) (Model0, error) {
	c := col()
	m := Model0{}
	err := FindDoc(c, eqCode(code), &m)
	return m, err
}

// Insert if a product does not already exist with the same code
func InsertModel0(m Model0) (int, error) {
	c := col()

	// Lock the collection for insert/update
	write.Lock()
	defer write.Unlock()
	return InsertDoc(c, eqCode(m.Code), encode(m))
}

func eqCode(code string) string {
	return `[{"eq": "` + code + `", "in": ["code"]}]`
}

func encode(m Model0) map[string]interface{} {
	return map[string]interface{}{"name": m.Name, "code": m.Code}
}

func col() *Col {
	return UseCol(MODEL0)
}

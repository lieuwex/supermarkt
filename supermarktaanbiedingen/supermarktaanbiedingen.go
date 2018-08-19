package supermarktaanbiedingen

import "supermarkt/supermarkts"

type AlbertHeijn struct{}

func (AlbertHeijn) Products(limit int) ([]supermarkts.Product, error) {
	return getProducts("ah", limit)
}
func (AlbertHeijn) ID() string {
	return "ah"
}
func (AlbertHeijn) Name() string {
	return "Albert Heijn"
}

type Dirk struct{}

func (Dirk) Products(limit int) ([]supermarkts.Product, error) {
	return getProducts("dirk", limit)
}
func (Dirk) ID() string {
	return "dirk"
}
func (Dirk) Name() string {
	return "Dirk"
}

type Hoogvliet struct{}

func (Hoogvliet) Products(limit int) ([]supermarkts.Product, error) {
	return getProducts("hoog", limit)
}
func (Hoogvliet) ID() string {
	return "hoog"
}
func (Hoogvliet) Name() string {
	return "Hoogvliet"
}

func init() {
	supermarkts.Register(&AlbertHeijn{})
	supermarkts.Register(&Dirk{})
	supermarkts.Register(&Hoogvliet{})
}

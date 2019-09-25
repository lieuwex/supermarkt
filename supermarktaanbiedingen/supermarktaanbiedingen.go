package supermarktaanbiedingen

import "supermarkt/supermarkts"

type SupermarktAanbiedingen struct{}

func (SupermarktAanbiedingen) Products(limit int) ([]supermarkts.Product, error) {
	return getProducts(limit)
}
func (SupermarktAanbiedingen) ID() string {
	return "supermarktaanbiedingen"
}
func (SupermarktAanbiedingen) Name() string {
	return "Supermarkt aanbiedingen"
}

func init() {
	supermarkts.Register(&SupermarktAanbiedingen{})
}

package supermarkts

import "fmt"

type SIUnit int

const (
	Unknown SIUnit = iota

	Liters
	Kilograms

	Unit
)

type Supermarket struct {
	ID   string
	Name string
}

type PriceInfo struct {
	Price float32

	Unit         SIUnit
	PricePerUnit float32
}

type Product struct {
	ID          string
	Name        string
	Brand       string
	Categories  []string
	Supermarket Supermarket

	PriceInfo PriceInfo
}

func (p Product) String() string {
	if p.PriceInfo.Unit == Unknown {
		return fmt.Sprintf("%s: [%s] %s (%v)", p.Supermarket.ID, p.ID, p.Name, p.PriceInfo.Price)
	}

	var unit string
	switch p.PriceInfo.Unit {
	case Liters:
		unit = "liter"
	case Kilograms:
		unit = "kilogram"

	case Unit:
		unit = "unit"
	}

	return fmt.Sprintf("%s: [%s] %s (%v/%s)", p.Supermarket.ID, p.ID, p.Name, p.PriceInfo.PricePerUnit, unit)
}

type Fetcher interface {
	ID() string
	Name() string

	Products(limit int) ([]Product, error)
}

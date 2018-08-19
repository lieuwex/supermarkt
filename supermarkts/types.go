package supermarkts

import "fmt"

type SIUnit int

const (
	Unknown SIUnit = iota

	Liters
	Kilograms

	Unit
)

type PriceInfo struct {
	Price float32

	Unit         SIUnit
	PricePerUnit float32
}

type Product struct {
	ID       string
	Name     string
	Brand    string
	Category string

	PriceInfo PriceInfo
}

func (p Product) String() string {
	if p.PriceInfo.Unit == Unknown {
		return fmt.Sprintf("[%s] %s (%v)", p.ID, p.Name, p.PriceInfo.Price)
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

	return fmt.Sprintf("[%s] %s (%v/%s)", p.ID, p.Name, p.PriceInfo.PricePerUnit, unit)
}

type Supermarket interface {
	ID() string
	Name() string

	Products(limit int) ([]Product, error)
}

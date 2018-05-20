package supermarkts

import (
	"errors"
	"fmt"
)

type SIUnit int

const (
	Liters SIUnit = iota
	Kilograms

	Unit
)

type PriceInfo struct {
	Cents int

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
	var unit string
	switch p.PriceInfo.Unit {
	case Liters:
		unit = "liters"
	case Kilograms:
		unit = "kilograms"

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

var drivers []Supermarket

func Register(sm Supermarket) {
	drivers = append(drivers, sm)
}

func getDriver(id string) Supermarket {
	for _, val := range drivers {
		if val.ID() == id {
			return val
		}
	}

	return nil
}

func ProductsBySupermarket(id string, limit int) ([]Product, error) {
	driver := getDriver(id)
	if driver == nil {
		return []Product{}, errors.New("no driver found")
	}

	return driver.Products(limit)
}

func Products(limit int) (map[string][]Product, error) {
	res := make(map[string][]Product)

	for _, driver := range drivers {
		items, err := driver.Products(limit)
		if err != nil {
			return res, err
		}

		res[driver.ID()] = items
	}

	return res, nil
}

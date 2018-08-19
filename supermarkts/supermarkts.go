package supermarkts

import "fmt"

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
		return []Product{}, fmt.Errorf("driver with ID %s not found", driver.ID())
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

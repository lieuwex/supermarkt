package main

import (
	"fmt"
	_ "supermarkt/jumbo"
	_ "supermarkt/supermarktaanbiedingen"
	"supermarkt/supermarkts"
)

func main() {
	items, err := supermarkts.ProductsBySupermarket("ah", 10)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Print("ah ")
		fmt.Println(item)
	}

	items, err = supermarkts.ProductsBySupermarket("dirk", 10)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Print("dirk ")
		fmt.Println(item)
	}
}

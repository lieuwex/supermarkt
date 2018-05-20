package main

import (
	"fmt"
	_ "supermarkt/jumbo"
	"supermarkt/supermarkts"
)

func main() {
	items, err := supermarkts.Products(100)
	if err != nil {
		panic(err)
	}

	for id, items := range items {
		fmt.Printf("\n----- %s -----\n", id)
		for _, val := range items {
			fmt.Println(val)
		}
	}
}

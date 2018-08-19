package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "supermarkt/jumbo"
	"supermarkt/supermarkts"
)

func main() {
	items, err := supermarkts.ProductsBySupermarket("jumbo", 1000)
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

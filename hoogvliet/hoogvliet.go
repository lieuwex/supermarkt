package main

import "supermarkt/supermarkts"

type Hoogvliet struct {
}

func init() {
	supermarkts.Register(&Hoogvliet{})
}

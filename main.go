package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	_ "supermarkt/jumbo"
	_ "supermarkt/supermarktaanbiedingen"
	"supermarkt/supermarkts"
	"time"
)

func main() {
	http.DefaultClient.Timeout = time.Hour

	m, err := supermarkts.Products(50000)
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("./output.json", bytes, 0644); err != nil {
		panic(err)
	}
}

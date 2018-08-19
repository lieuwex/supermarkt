package jumbo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"supermarkt/supermarkts"

	"github.com/PuerkitoBio/goquery"
)

type Jumbo struct{}

func (Jumbo) ID() string   { return "jumbo" }
func (Jumbo) Name() string { return "Jumbo" }

func init() {
	supermarkts.Register(&Jumbo{})
}

var siRegexp = regexp.MustCompile(`\(([\d,]+)/(\w+)\)`)

func parseSIPrice(str string) (price float32, unit supermarkts.SIUnit, ok bool) {
	res := siRegexp.FindStringSubmatch(str)
	if res == nil {
		return 0, 0, false
	}

	priceStr := strings.Replace(res[1], ",", ".", -1)
	priceEuros, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		return 0, 0, false
	}

	switch res[2] {
	case "kilo":
		unit = supermarkts.Kilograms
	case "liter":
		unit = supermarkts.Liters

	case "stuk":
		unit = supermarkts.Unit
	}

	return float32(priceEuros), unit, true
}

func (Jumbo) getPage(n int) ([]supermarkts.Product, error) {
	url := fmt.Sprintf("https://www.jumbo.com/producten?PageNumber=%d", n)
	res, err := http.Get(url)
	if err != nil {
		return []supermarkts.Product{}, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return []supermarkts.Product{}, err
	}

	var products []supermarkts.Product
	doc.Find(".jum-item-product").Each(func(n int, s *goquery.Selection) {
		val, ok := s.Attr("data-jum-product-impression")
		if !ok {
			return
		}

		var m map[string]interface{}
		if err := json.Unmarshal([]byte(val), &m); err != nil {
			return
		}

		brand, ok := s.Attr("data-jum-brand")
		if !ok {
			return
		}

		priceElem := s.Find("#PriceInCents_" + m["id"].(string)).First()
		if priceElem == nil {
			return
		}
		centsStr, ok := priceElem.Attr("value")
		if !ok {
			return
		}
		cents, err := strconv.Atoi(centsStr)
		if err != nil {
			return
		}

		price, unit, ok := parseSIPrice(s.Find(".jum-comparative-price").Text())

		product := supermarkts.Product{
			ID:       m["id"].(string),
			Name:     m["name"].(string),
			Brand:    brand,
			Category: m["category"].(string),
			PriceInfo: supermarkts.PriceInfo{
				Price:        float32(cents) / 100,
				Unit:         unit,
				PricePerUnit: price,
			},
		}
		products = append(products, product)
	})

	return products, nil
}

func (j Jumbo) Products(limit int) ([]supermarkts.Product, error) {
	var res []supermarkts.Product

	for i := 0; ; i++ {
		items, err := j.getPage(i)
		if err != nil {
			return []supermarkts.Product{}, err
		} else if len(items) == 0 {
			break
		}

		res = append(res, items...)

		if len(res) >= limit {
			break
		}
	}

	if len(res) >= limit {
		res = res[:limit]
	}

	return res, nil
}

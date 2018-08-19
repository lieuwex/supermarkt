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
	"golang.org/x/net/html"
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

func getPage(n int) ([]supermarkts.Product, error) {
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

	nodes := doc.Find(".jum-item-product").Nodes
	var products []supermarkts.Product
	for _, node := range nodes {
		s := goquery.Selection{
			Nodes: []*html.Node{node},
		}

		val, ok := s.Attr("data-jum-product-impression")
		if !ok {
			continue
		}
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(val), &m); err != nil {
			continue
		}
		id := m["id"].(string)

		brand, ok := s.Attr("data-jum-brand")
		if !ok {
			continue
		}

		priceElem := s.Find("#PriceInCents_" + id).First()
		if priceElem == nil {
			continue
		}
		centsStr, ok := priceElem.Attr("value")
		if !ok {
			continue
		}
		cents, err := strconv.Atoi(centsStr)
		if err != nil {
			continue
		}

		siPrice := s.Find(".jum-comparative-price").Text()
		pricePerUnit, unit, ok := parseSIPrice(siPrice)

		product := supermarkts.Product{
			ID:       id,
			Name:     m["name"].(string),
			Brand:    brand,
			Category: m["category"].(string),
			PriceInfo: supermarkts.PriceInfo{
				Price:        float32(cents) / 100,
				Unit:         unit,
				PricePerUnit: pricePerUnit,
			},
		}
		products = append(products, product)
	}

	return products, nil
}

func (Jumbo) Products(limit int) ([]supermarkts.Product, error) {
	var res []supermarkts.Product

	for i := 0; ; i++ {
		items, err := getPage(i)
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

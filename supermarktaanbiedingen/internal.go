package supermarktaanbiedingen

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"supermarkt/supermarkts"

	"github.com/Jeffail/tunny"
	"github.com/PuerkitoBio/goquery"
	"github.com/namsral/microdata"
	"golang.org/x/net/html"
)

type Offer struct {
	Price       float32
	Supermarket string
}

type Product struct {
	Name   string
	Brand  string
	Offers []Offer
}

func supermarketToID(name string) string {
	return map[string]string{
		"Albert Heijn XL": "ah",
		"Albert Heijn":    "ah",
		"Boni":            "boni",
		"Coop":            "coop",
		"Deen":            "deen",
		"Dekamarkt":       "deka",
		"Dirk":            "dirk",
		"EMTE":            "emte",
		"Hoogvliet":       "hoog",
		"Jan Linders":     "jan",
		"Lidl":            "lidl",
		"Nettorama":       "nett",
		"Plus":            "plus",
		"Poiesz":          "poiesz",
		"Spar":            "spar",
		"Vomar":           "vomar",
	}[name]
}

func getProduct(url string) (Product, error) {
	data, err := microdata.ParseURL(url)
	if err != nil {
		return Product{}, err
	}

	prod := data.Items[0].Properties

	brandRaw, ok := prod["brand"]
	brand := ""
	if ok {
		brand = brandRaw[0].(*microdata.Item).Properties["name"][0].(string)
	}

	var offers []Offer
	for _, x := range prod["offers"] {
		offer := x.(*microdata.Item).Properties
		supermarketRaw, ok := offer["availableAtOrFrom"]
		supermarket := ""
		if ok {
			supermarket = supermarketRaw[0].(*microdata.Item).Properties["name"][0].(string)
		}

		priceRaw, ok := offer["price"]
		priceStr := ""
		if ok {
			priceStr = priceRaw[0].(string)
		}

		v, _ := strconv.ParseFloat(strings.Replace(priceStr, ",", ".", -1), 32)
		offers = append(offers, Offer{
			Price:       float32(v),
			Supermarket: supermarket,
		})
	}

	return Product{
		Name:   strings.TrimSpace(prod["name"][0].(string)),
		Brand:  brand,
		Offers: offers,
	}, nil
}

type workerRes struct {
	Result Product
	Error  error
}

func handleResultsPage(n int) ([]Product, bool, error) {
	// download search page

	url := fmt.Sprintf(
		"http://www.supermarktaanbiedingen.com/zoeken/%%20/pagina/%d",
		n,
	)
	res, err := http.Get(url)
	if err != nil {
		return []Product{}, false, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return []Product{}, false, err
	}

	// collect product urls

	nodes := doc.Find(".price, .card_prijs-oud").Nodes
	any := len(nodes) > 0
	var urls []string
	for _, node := range nodes {
		s := goquery.Selection{
			Nodes: []*html.Node{node},
		}

		if s.HasClass("card_prijs-oud") {
			continue
		}

		a := s.ParentsFiltered("a")
		slug, ok := a.Attr("href")
		if !ok {
			continue
		}

		url := "http://" + path.Join("www.supermarktaanbiedingen.com/", slug)
		urls = append(urls, url)
	}

	// scrape product pages

	pool := tunny.NewFunc(5, func(url interface{}) interface{} {
		product, err := getProduct(url.(string))
		return workerRes{
			Result: product,
			Error:  err,
		}
	})

	var products []Product
	for _, url := range urls {
		res := pool.Process(url).(workerRes)
		if res.Error != nil {
			return []Product{}, false, res.Error
		}

		products = append(products, res.Result)
	}
	return products, any, nil
}

func getProducts(limit int) ([]supermarkts.Product, error) {
	var res []supermarkts.Product

	for i := 1; ; i++ {
		prods, contains, err := handleResultsPage(i)
		if err != nil {
			panic(err)
		}

		for _, p := range prods {
			for _, offer := range p.Offers {
				p := supermarkts.Product{
					ID:         "", // TODO
					Name:       p.Name,
					Brand:      p.Brand,
					Categories: []string{}, // TODO
					Supermarket: supermarkts.Supermarket{
						ID:   supermarketToID(offer.Supermarket),
						Name: offer.Supermarket,
					},
					PriceInfo: supermarkts.PriceInfo{
						Price: offer.Price,
					},
				}
				res = append(res, p)
			}
		}

		if !contains || len(res) >= limit {
			break
		}
	}

	return res, nil
}

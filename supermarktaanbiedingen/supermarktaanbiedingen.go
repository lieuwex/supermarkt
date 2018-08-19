package supermarktaanbiedingen

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"supermarkt/supermarkts"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/namsral/microdata"
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

func handleResultsPage(n int) ([]Product, bool, error) {
	url := fmt.Sprintf("http://www.supermarktaanbiedingen.com/zoeken/%%20/pagina/%d", n)
	resp, err := http.Get(url)
	if err != nil {
		return []Product{}, false, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []Product{}, false, err
	}

	prices := doc.Find(".price, .card_prijs-oud")
	urlsCh := make(chan string, prices.Length())

	if prices.Length() == 0 {
		return []Product{}, false, nil
	}

	prices.Each(func(n int, s *goquery.Selection) {
		if s.HasClass("card_prijs-oud") {
			return
		}

		a := s.Parent().Parent().Parent().Parent() // HACK
		slug, ok := a.Attr("href")
		if !ok {
			return
		}

		urlsCh <- "http://" + path.Join("www.supermarktaanbiedingen.com/", slug)
	})

	nURLs := len(urlsCh)
	var wg sync.WaitGroup
	wg.Add(nURLs)
	products := make(chan Product, nURLs)

	for i := 0; i < nURLs; i++ {
		go func() {
			url := <-urlsCh
			product, err := getProduct(url)
			if err == nil {
				products <- product
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(urlsCh)
	close(products)

	var res []Product
	for prod := range products {
		res = append(res, prod)
	}
	return res, true, nil
}

type AlbertHeijn struct{}

func (AlbertHeijn) Products(limit int) ([]supermarkts.Product, error) {
	var res []supermarkts.Product

	for i := 1; ; i++ {
		prods, contains, err := handleResultsPage(i)
		if err != nil {
			panic(err)
		}

		for _, p := range prods {
			fmt.Printf("%#v\n", p)
		}

		if !contains {
			break
		}
	}

	return res, nil
}

func (AlbertHeijn) ID() string   { return "ah" }
func (AlbertHeijn) Name() string { return "Albert Heijn" }

func init() {
	supermarkts.Register(&AlbertHeijn{})
}

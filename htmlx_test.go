package htmlx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type Categories struct {
	TravelPageUrl  string `htmlx_sel:"ul > li > ul > li:nth-child(1) > a" htmlx_src:"attr(href)"`
	MysteryPageUrl string `htmlx_sel:"ul > li > ul > li:nth-child(2) > a" htmlx_src:"attr(href)"`
}

type TestStruct struct {
	SiteName        string         `htmlx_sel:"#default > header > div > div > div > a"`
	NumberOfResults *int           `htmlx_sel:"#default > div > div > div > div > form > strong:nth-child(2)"`
	SiteUrl         string         `htmlx_sel:"#default > header > div > div > div > a"                                                                                            htmlx_src:"attr(href)"`
	TopBookPrice    FloatUnitValue `htmlx_sel:"#default > div > div > div > div > section > div:nth-child(2) > ol > li:nth-child(1) > article > div.product_price > p.price_color"`
	Categories      Categories     `htmlx_sel:"#default > div > div > div > aside > div.side_categories"                                                                           htmlx_src:"_"`
}

func TestHtmlxNode(t *testing.T) {
	res, err := http.Get("https://books.toscrape.com")
	if err != nil {
		t.Fatal(err)
	}

	for range 100 {
		num := 0

		testStruct := TestStruct{NumberOfResults: &num}

		htmlxDoc, err := NewDocFromReader(res.Body, Async(true))
		if err != nil {
			t.Fatal(err)
		}

		htmlxDoc.Scan(&testStruct)

		jsonDat, err := json.MarshalIndent(testStruct, "", "	")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(jsonDat))
	}
}

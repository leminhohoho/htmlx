package htmlx

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type Categories struct {
	TravelPageUrl  string `htmlx_sel:"ul > li > ul > li:nth-child(1) > a" htmlx_src:"attr(href)"`
	MysteryPageUrl string `htmlx_sel:"ul > li > ul > li:nth-child(2) > a" htmlx_src:"attr(href)"`
}

type TestStruct struct {
	SiteName   string     `htmlx_sel:"#default > header > div > div > div > a"`
	SiteUrl    string     `htmlx_sel:"#default > header > div > div > div > a"                  htmlx_src:"attr(href)"`
	Categories Categories `htmlx_sel:"#default > div > div > div > aside > div.side_categories" htmlx_src:"_"`
}

func TestHtmlxNode(t *testing.T) {
	res, err := http.Get("https://books.toscrape.com")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	for range 10000 {
		testStruct := &TestStruct{}

		rootNode := &HtmlxNode{
			Selection: doc.Selection,
			val:       reflect.ValueOf(*testStruct),
			config:    &Config{async: true},
		}

		if err := rootNode.Construct(); err != nil {
			t.Fatal(err)
		}

		fmt.Printf("%+v\n", rootNode)
	}
}

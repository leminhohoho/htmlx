package htmlx

import (
	"encoding/json"
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
	SiteName        string     `htmlx_sel:"#default > header > div > div > div > a"`
	NumberOfResults *int       `htmlx_sel:"#default > div > div > div > div > form > strong:nth-child(2)"`
	SiteUrl         string     `htmlx_sel:"#default > header > div > div > div > a"                       htmlx_src:"attr(href)"`
	Categories      Categories `htmlx_sel:"#default > div > div > div > aside > div.side_categories"      htmlx_src:"_"`
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

	for range 100 {
		num := 0

		testStruct := &TestStruct{NumberOfResults: &num}

		rootNode := &HtmlxNode{
			Selection: doc.Selection,
			val:       reflect.ValueOf(testStruct).Elem(),
			config:    &Config{async: true},
		}

		if err := rootNode.Construct(); err != nil {
			t.Fatal(err)
		}

		if err := rootNode.Parse(); err != nil {
			t.Fatal(err)
		}

		jsonDat, err := json.MarshalIndent(*testStruct, "", "	")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(jsonDat))
	}
}

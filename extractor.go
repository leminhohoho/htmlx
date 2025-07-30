package htmlx

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Extractor take a [github.com/PuerkitoBio/goquery.Selection] and return a string value extracted from the selection.
// If the extractor need to use external data, consider using higher order functions with closures.
type Extractor func(*goquery.Selection) (string, error)

func extractText(s *goquery.Selection) (string, error) {
	return s.Clone().Children().Remove().End().Text(), nil
}

func extractHtml(s *goquery.Selection) (string, error) {
	return s.Html()
}

func extractAttr(attrName string) Extractor {
	return func(s *goquery.Selection) (string, error) {
		v, exists := s.Attr(attrName)
		if !exists {
			return "", fmt.Errorf("Unable to locate attr '%s'", attrName)
		}

		return v, nil
	}
}

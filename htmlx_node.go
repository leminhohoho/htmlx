package htmlx

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Htmlx node represent a collection of HTML nodes which is binded to a struct field which receive the data.
// Htmlx node using chainable j-query like syntax, which call the according method from [github.com/PuerkitoBio/goquery.Selection] under the hood.
// Htmlx does not replace [github.com/PuerkitoBio/goquery.Selection], so it will only expose common methods like Find(), Remove(),...
// To take full advantage of [github.com/PuerkitoBio/goquery], accessing [HtmlxNode.Selection]
type HtmlxNode struct {
	// contains filtered or unexported fields

	mu sync.Mutex

	Selection *goquery.Selection
	name      string
	config    *Config
	extractor Extractor
	val       reflect.Value
	children  []*HtmlxNode

	constructed bool
}

func getExtractor(src string) (Extractor, error) {
	if src == "_" {
		return nil, nil
	} else if src == "text" || src == "" {
		return extractText, nil
	} else if src == "html" {
		return extractHtml, nil
	} else if regexp.MustCompile(`^attr\([a-zA-Z-_:]+\)$`).MatchString(src) {
		attrName := src[5 : len(src)-1]
		return extractAttr(attrName), nil
	}

	return nil, fmt.Errorf("'%s' is not a valid src", src)
}

func (n *HtmlxNode) appendNode(node *HtmlxNode) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.children = append(n.children, node)
}

func (n *HtmlxNode) registerNode(fieldVal reflect.Value, fieldInfo reflect.StructField) error {
	selector := fieldInfo.Tag.Get("htmlx_sel")
	if selector == "" {
		return nil
	}

	src := strings.TrimSpace(fieldInfo.Tag.Get("htmlx_src"))

	extractor, err := getExtractor(src)
	if err != nil {
		return err
	}

	node := HtmlxNode{
		Selection: n.Selection.Find(selector),
		name:      fieldInfo.Name,
		config:    n.config,
		extractor: extractor,
		val:       fieldVal,
	}

	if err := node.Construct(); err != nil {
		return err
	}

	n.appendNode(&node)

	return nil
}

// Construct create a [HtmlxNode] tree from the current root node.
// Construct can only be called one provided that it is successful, subsequent calls will return error
func (n *HtmlxNode) Construct() error {
	if n.constructed {
		return fmt.Errorf("The node is already constructed")
	}

	if n.val.Kind() != reflect.Struct {
		return nil
	}

	var wg sync.WaitGroup

	for i := range n.val.NumField() {
		fieldVal := n.val.Field(i)
		fieldInfo := n.val.Type().Field(i)

		var err error

		if n.config.async {
			wg.Add(1)
			go func() { err = n.registerNode(fieldVal, fieldInfo); wg.Done() }()
		} else {
			err = n.registerNode(fieldVal, fieldInfo)
		}
		if err != nil {
			n.children = nil
			return &ErrConstructHtmlxNode{fieldInfo.Name, err}
		}
	}

	wg.Wait()
	n.constructed = true

	return nil
}

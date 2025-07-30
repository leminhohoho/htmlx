package htmlx

import (
	"encoding"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Htmlx node represent a collection of HTML nodes which is binded to a struct field which receive the data.
type htmlxNode struct {
	// contains filtered or unexported fields

	mu sync.Mutex

	selection *goquery.Selection // Contain the HTML node to extract from
	name      string
	config    *Config
	extractor Extractor
	val       reflect.Value
	children  []*htmlxNode

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

func (n *htmlxNode) appendNode(node *htmlxNode) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.children = append(n.children, node)
}

func (n *htmlxNode) registerNode(fieldVal reflect.Value, fieldInfo reflect.StructField) error {
	selector := fieldInfo.Tag.Get("htmlx_sel")
	if selector == "" {
		return nil
	}

	src := strings.TrimSpace(fieldInfo.Tag.Get("htmlx_src"))

	extractor, err := getExtractor(src)
	if err != nil {
		return err
	}

	node := htmlxNode{
		selection: n.selection.Find(selector),
		name:      fieldInfo.Name,
		config:    n.config,
		extractor: extractor,
		val:       fieldVal,
	}

	if err := node.construct(); err != nil {
		return err
	}

	n.appendNode(&node)

	return nil
}

// construct create a [htmlxNode] tree from the current root node.
// construct can only be called once provided that it is successful, calls after that will return error.
func (n *htmlxNode) construct() error {
	if n.constructed {
		return fmt.Errorf("The node is already constructed")
	}

	if n.val.Kind() != reflect.Struct {
		n.constructed = true
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

func (n *htmlxNode) parseFromSelf() error {
	if n.extractor == nil {
		return nil
	}

	rawVal, err := n.extractor(n.selection)
	if err != nil {
		return err
	}

	ptr := reflect.New(n.val.Type())
	if ptr.Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
		marshaller, _ := ptr.Interface().(encoding.TextUnmarshaler)
		if err := marshaller.UnmarshalText([]byte(rawVal)); err != nil {
			return err
		}

		n.val.Set(reflect.ValueOf(marshaller).Elem())
	}

	switch n.val.Kind() {
	case reflect.Ptr:
		if n.val.IsNil() {
			n.val.Set(reflect.New(n.val.Type().Elem()))
		}

		n.val = n.val.Elem()
		return n.parseFromSelf()
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		num, err := strconv.Atoi(rawVal)
		if err != nil {
			return err
		}

		n.val.SetInt(int64(num))
	case reflect.Float32, reflect.Float64:
		num, err := strconv.ParseFloat(rawVal, 64)
		if err != nil {
			return err
		}

		n.val.SetFloat(num)
	case reflect.String:
		n.val.SetString(rawVal)
	case reflect.Uint8:
		n.val.SetBytes([]byte(rawVal))
	default:
		return fmt.Errorf("Value of type '%v' is not supported", n.val.Type())
	}

	return nil
}

// parse extract from HTMl content and parse onto struct.
// parse will return error if node is not constructed.
func (n *htmlxNode) parse() error {
	if !n.constructed {
		return &ErrParseHtmlxNode{n.name, fmt.Errorf("Htmlx node is not constructed")}
	}

	var err error

	if err = n.parseFromSelf(); err != nil {
		return &ErrParseHtmlxNode{n.name, err}
	}

	var wg sync.WaitGroup

	for _, child := range n.children {
		if n.config.async {
			wg.Add(1)
			go func() { err = child.parse(); wg.Done() }()
		} else {
			err = child.parse()
		}
		if err != nil {
			return &ErrParseHtmlxNode{n.name, err}
		}
	}

	wg.Wait()

	return nil
}

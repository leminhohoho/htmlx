package htmlx

import (
	"fmt"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

// Selection is a wrapper of [github.com/PuerkitoBio/goquery.Selection] with additional functionality.
type Selection struct {
	// contains filtered or unexported fields

	*goquery.Selection
	config Config
}

// Clonex is a wrapper of [github.com/PuerkitoBio/goquery.Selection.Clone] with additional functionality.
// It creates a deep copy of the set of matched nodes.
// The new nodes will not be attached to the document.
func (s *Selection) Clonex(opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	return &Selection{Selection: s.Clone(), config: c}
}

// Findx is a wrapper of [github.com/PuerkitoBio/goquery.Selection.Find] with additional functionality.
// It gets the descendants of each element in the current set of matched elements, filtered by a selector.
// It returns a new Selection object containing these matched elements.
func (s *Selection) Findx(selector string, opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	return &Selection{Selection: s.Find(selector), config: c}
}

// Childrenx is a wrapper of [github.com/PuerkitoBio/goquery.Selection.Children] with additional functionality.
// It gets the child elements of each element in the Selection.
// It returns a new Selection object containing these elements.
func (s *Selection) Childrenx(opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	return &Selection{Selection: s.Children(), config: c}
}

// Firstx is a wrapper of [github.com/PuerkitoBio/goquery.Selection.First] with additional functionality.
// It reduces the set of matched elements to the first in the set.
// It returns a new Selection object, and an empty Selection object if the the selection is empty.
func (s *Selection) Firstx(opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	return &Selection{Selection: s.First(), config: c}
}

// Endx is a wrapper of [github.com/PuerkitoBio/goquery.Selection.End] with additional functionality.
// It ends the most recent filtering operation in the current chain and returns the set of matched elements to its previous state.
func (s *Selection) Endx(opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	return &Selection{Selection: s.End(), config: c}
}

// Removex is a wrapper of [github.com/PuerkitoBio/goquery.Selection.Remove] with additional functionality.
// It removes the set of matched elements from the document.
// It returns the same selection, now consisting of nodes not in the document.
func (s *Selection) Removex(opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	return &Selection{Selection: s.Remove(), config: c}
}

// Eachx is a wrapper of [github.com/PuerkitoBio/goquery.Selection.Each] with additional functionality.
// It iterate through each matched element in the selection object and called f for each one.
// The index start at 0.
// It return the current selection.
func (s *Selection) Eachx(f func(int, *Selection), opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	s.Each(func(i int, selection *goquery.Selection) { f(i, &Selection{Selection: selection, config: c}) })

	return s
}

// Eqx is a wrapper of [github.com/PuerkitoBio/goquery.Selection.Eq] with additional functionality.
// It reduces the set of matched elements to the one at the specified index.
// If a negative index is given, it counts backwards starting at the end of the set.
// It returns a new Selection object, and an empty Selection object if the index is invalid.
func (s *Selection) Eqx(index int, opts ...Option) *Selection {
	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	return &Selection{Selection: s.Eq(index), config: c}
}

// Scan extract the data from HTMl content into the value pointed at by dest, if dest is not a pointer, the method will return error.
// opts are options which will be appended to the current config of the selection.
// It return the current selection and 2nd value as error.
func (s *Selection) Scan(dest any, opts ...Option) (*Selection, error) {
	p := reflect.ValueOf(dest)

	if p.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Error: '%v' is not a pointer", p.Type())
	}

	if p.IsNil() {
		return nil, fmt.Errorf("Error: dest is a nil pointer")
	}

	c := s.config
	for _, opt := range opts {
		opt(&c)
	}

	rootNode := &htmlxNode{
		selection: s.Selection,
		val:       p.Elem(),
		config:    &c,
	}

	if err := rootNode.construct(); err != nil {
		return nil, err
	}

	if err := rootNode.parse(); err != nil {
		return nil, err
	}

	return s, nil
}

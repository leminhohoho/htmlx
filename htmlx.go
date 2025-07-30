package htmlx

import (
	"bytes"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Document struct {
	// contains filtered or unexported fields

	*Selection
}

// NewDocFromReader return a Document from an [io.Reader].
// It does not check if the reader is also an io.Closer, the provided reader is never closed by this call.
// It is the responsibility of the caller to close it if required.
// It return a 2nd value as error if the reader's data can't be parsed as html.
// opts are options which will be appended to the current config of the selection.
func NewDocFromReader(r io.Reader, opts ...Option) (*Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	c := Config{}
	for _, opt := range opts {
		opt(&c)
	}

	return &Document{Selection: &Selection{Selection: doc.Selection, config: c}}, nil
}

// NewDocFromString create an [io.Reader] from given string and call [NewDocFromReader] under the hood.
func NewDocFromString(str string, opts ...Option) (*Document, error) {
	return NewDocFromReader(strings.NewReader(str), opts...)
}

// NewDocFromBytes create an [io.Reader] from given bytes and call [NewDocFromReader] under the hood.
func NewDocFromBytes(b []byte, opts ...Option) (*Document, error) {
	return NewDocFromReader(bytes.NewReader(b), opts...)
}

package htmlx

import "fmt"

type ErrConstructHtmlxNode struct {
	fieldName string
	err       error
}

func (e *ErrConstructHtmlxNode) Error() string {
	return fmt.Sprintf("Error constructing field '%s': %v", e.fieldName, e.err)
}

type ErrParseHtmlxNode struct {
	fieldName string
	err       error
}

func (e *ErrParseHtmlxNode) Error() string {
	return fmt.Sprintf("Error parsing field '%s': %v", e.fieldName, e.err)
}

package htmlx

import "fmt"

type ErrConstructHtmlxNode struct {
	fieldName string
	err       error
}

func (e *ErrConstructHtmlxNode) Error() string {
	return fmt.Sprintf("Error processing field '%s': %v", e.fieldName, e.err)
}

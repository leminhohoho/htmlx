package htmlx

import (
	"fmt"
	"regexp"
	"strconv"
)

// FloatUnitValue is a wrapper for float64 which implement [encoding.TextUnmarshaler] to convert number with unit to float64.
type FloatUnitValue float64

func (n *FloatUnitValue) UnmarshalText(text []byte) error {
	numStr := regexp.MustCompile(`-?\d+(?:[,.]\d+)*(\.\d+)?`).FindString(string(text))
	if numStr == "" {
		return fmt.Errorf("Unable to extract number from '%s'", string(text))
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return err
	}

	*n = FloatUnitValue(num)

	return nil
}

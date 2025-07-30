package htmlx

import (
	"fmt"
	"regexp"
	"strconv"
)

// FloatUnitValue is a wrapper for float64 which implement [encoding.TextUnmarshaler] to convert number with unit to float64.
type FloatUnitValue float64

// UnmarshalText implement the [encoding.TextUnmarshaler] interface
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

// IntUnitValue is a wrapper for int which implement [encoding.TextUnmarshaler] to convert number with unit to int.
type IntUnitValue int

// UnmarshalText implement the [encoding.TextUnmarshaler] interface
func (n *IntUnitValue) UnmarshalText(text []byte) error {
	floatUnitValue := FloatUnitValue(*n)
	if err := floatUnitValue.UnmarshalText(text); err != nil {
		return err
	}

	*n = IntUnitValue(floatUnitValue)

	return nil
}

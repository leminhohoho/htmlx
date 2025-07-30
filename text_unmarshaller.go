package htmlx

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// FloatUnitValue is a wrapper for float64 which implement [encoding.TextUnmarshaler] to convert number with unit to float64.
type FloatUnitValue float64

// UnmarshalText implement the [encoding.TextUnmarshaler] interface
func (n *FloatUnitValue) UnmarshalText(text []byte) error {
	numStr := regexp.MustCompile(`-?\d+(?:[,.]\d+)*(\.\d+)?`).FindString(string(text))
	if numStr == "" {
		return fmt.Errorf("Unable to extract number from '%s'", string(text))
	}
	numStr = strings.ReplaceAll(numStr, ",", "")

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

// Time is a wrapper for time.Time which implement [encoding.TextUnmarshaler] to convert string to date.
// Time use [Time.Layout] to determine the date layout which is used by [time.Parse], if not given it is default to 2006-01-02T15:04:05Z07:00.
// After converting, the value can be extracted from [Time.Time]
type Time struct {
	Layout string
	Time   time.Time
}

// UnmarshalText implement the [encoding.TextUnmarshaler] interface
func (t *Time) UnmarshalText(text []byte) error {
	var err error

	if t.Layout == "" {
		t.Layout = "2006-01-02T15:04:05Z07:00"
	}

	t.Time, err = time.Parse(t.Layout, string(text))
	if err != nil {
		return err
	}

	return nil
}

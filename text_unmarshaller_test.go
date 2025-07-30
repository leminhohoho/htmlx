package htmlx

import (
	"testing"
	"time"
)

func TestFloatUnitValue(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{"20k", 20.0},
		{"$59.99", 59.99},
		{"1,000,000", 1000000.0},
		{"-10.5", -10.5},
		{"100 USD", 100.0},
		{"0.123", 0.123},
		{"-0.5", -0.5},
		{"1,000,000.00", 1000000.0},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var num FloatUnitValue
			err := num.UnmarshalText([]byte(tc.input))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if float64(num) != tc.expected {
				t.Errorf("Expected %f, but got %f", tc.expected, float64(num))
			}
		})
	}
}

func TestIntUnitValue(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"20k", 20},
		{"$59.99", 59},
		{"1,000,000", 1000000},
		{"-10.5", -10},
		{"100 USD", 100},
		{"-0.5", 0},
		{"1,000,000.00", 1000000},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var num IntUnitValue
			err := num.UnmarshalText([]byte(tc.input))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if int(num) != tc.expected {
				t.Errorf("Expected %d, but got %d", tc.expected, int(num))
			}
		})
	}
}

func TestTimeValue(t *testing.T) {
	testCases := []struct {
		input    string
		layout   string
		expected time.Time
	}{
		{"2024-01-15T10:30:00Z", "2006-01-02T15:04:05Z", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)},
		{"31-07-2024", "02-01-2006", time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)},
		{"2023-Mar-10 15:04:05", "2006-Jan-02 15:04:05", time.Date(2023, 3, 10, 15, 4, 5, 0, time.UTC)},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var tv Time
			tv.Layout = tc.layout
			err := tv.UnmarshalText([]byte(tc.input))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tv.Time.Equal(tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, tv.Time)
			}
		})
	}
}

func TestFloatUnitValueError(t *testing.T) {
	var num FloatUnitValue
	err := num.UnmarshalText([]byte("abc"))
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestIntUnitValueError(t *testing.T) {
	var num IntUnitValue
	err := num.UnmarshalText([]byte("abc"))
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestTimeValueError(t *testing.T) {
	var tv Time
	err := tv.UnmarshalText([]byte("invalid-time"))
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}


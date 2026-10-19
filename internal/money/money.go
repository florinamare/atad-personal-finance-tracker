// Package money represents amounts exactly, as integer cents.
// An amount never passes through a floating-point number.
package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Amount is a sum of money in cents: "-12.50" is -1250.
type Amount int64

var ErrInvalid = errors.New("invalid amount")

// At most 13 digits before the point: far from int64 overflow even when summed.
const maxIntDigits = 13

// Parse accepts exactly the API format: optional "-", digits, ".", two digits.
func Parse(s string) (Amount, error) {
	orig := s
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}

	intPart, fracPart, ok := strings.Cut(s, ".")
	if !ok || intPart == "" || len(fracPart) != 2 ||
		len(intPart) > maxIntDigits || !allDigits(intPart) || !allDigits(fracPart) {
		return 0, fmt.Errorf("%w: %q (expected a string like \"-12.50\")", ErrInvalid, orig)
	}

	i, _ := strconv.ParseInt(intPart, 10, 64)
	f, _ := strconv.ParseInt(fracPart, 10, 64)
	cents := i*100 + f
	if neg {
		cents = -cents
	}
	return Amount(cents), nil
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// String formats back: -1250 -> "-12.50".
func (a Amount) String() string {
	v := int64(a)
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	return fmt.Sprintf("%s%d.%02d", sign, v/100, v%100)
}

// MarshalJSON writes the amount as a JSON string, never a number.
func (a Amount) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON accepts only a JSON string such as "-12.50".
func (a *Amount) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("%w: amount must be a JSON string like \"-12.50\", not a number", ErrInvalid)
	}
	v, err := Parse(s)
	if err != nil {
		return err
	}
	*a = v
	return nil
}
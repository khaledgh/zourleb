// Package phone normalizes and validates Lebanese phone numbers to E.164.
package phone

import (
	"errors"
	"regexp"
	"strings"
)

// ErrInvalid is returned when a number cannot be normalized to a valid
// Lebanese E.164 number.
var ErrInvalid = errors.New("invalid lebanese phone number")

const countryCode = "961"

// Lebanese subscriber numbers are 7 or 8 digits after the country code.
// Mobile: 3/7/8/76/78/79/81 prefixes; landline: 1/4/5/6/7/8/9.
var localPattern = regexp.MustCompile(`^[0-9]{7,8}$`)

var nonDigit = regexp.MustCompile(`[^0-9]`)

// Normalize converts a variety of input formats to E.164 (+961XXXXXXXX).
// Accepts: "+961 3 123456", "0096170123456", "03 123 456", "70123456".
func Normalize(raw string) (string, error) {
	if raw == "" {
		return "", ErrInvalid
	}
	s := strings.TrimSpace(raw)
	hadPlus := strings.HasPrefix(s, "+")
	digits := nonDigit.ReplaceAllString(s, "")

	switch {
	case strings.HasPrefix(digits, "00"+countryCode):
		digits = digits[len("00"+countryCode):]
	case hadPlus && strings.HasPrefix(digits, countryCode):
		digits = digits[len(countryCode):]
	case strings.HasPrefix(digits, countryCode) && len(digits) > 8:
		// Bare country code without plus, e.g. "96170123456".
		digits = digits[len(countryCode):]
	case strings.HasPrefix(digits, "0"):
		// National trunk prefix, e.g. "03123456".
		digits = strings.TrimPrefix(digits, "0")
	}

	if !localPattern.MatchString(digits) {
		return "", ErrInvalid
	}
	return "+" + countryCode + digits, nil
}

// IsValid reports whether raw normalizes to a valid Lebanese number.
func IsValid(raw string) bool {
	_, err := Normalize(raw)
	return err == nil
}

// Mask hides the middle of a normalized number for logs/UI, e.g. +96170***456.
func Mask(e164 string) string {
	if len(e164) < 7 {
		return "***"
	}
	return e164[:6] + "***" + e164[len(e164)-3:]
}

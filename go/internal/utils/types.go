package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// UUIDsToStrings converts a slice of uuid.UUID objects to a slice of their string representations.
//
// Parameters:
//   - uuids: A slice of uuid.UUID objects to be converted.
//
// Returns:
//   - []string: A slice containing the string representations of the input UUIDs.
func UUIDsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, guid := range uuids {
		strs[i] = guid.String()
	}
	return strs
}

// ExtractAmountParts splits a monetary amount into its integer (left) and fractional (right) parts.
// It supports input of type string, uint32, and uint64.
//
// For string inputs, the function expects a decimal representation of the amount
// (e.g., "123.45"). If the fractional part is missing, it defaults to "00". If the
// fractional part has only one digit, it is padded with a trailing zero.
//
// For uint32 and uint64 inputs, the function assumes the amount is represented
// in cents (e.g., 12345 represents "123.45"). The integer part is derived by dividing
// the amount by 100, and the fractional part is derived from the remainder. If the
// fractional part has only one digit, it is padded with a leading zero.
//
// Parameters:
// - amount: The monetary amount to be split. It can be of type string, uint32, or uint64.
//
// Returns:
// - left: The integer part of the amount as a string.
// - rignt: The fractional part of the amount as a string (always two digits).
func ExtractAmountParts(amount any) (left string, rignt string) {

	switch amount := amount.(type) {
	case string:
		splitedAmount := strings.Split(amount, ".")
		left = splitedAmount[0]
		if len(splitedAmount) == 1 {
			rignt = "00"
		} else if len(splitedAmount[1]) == 1 {
			rignt = splitedAmount[1] + "0"
		} else {
			rignt = splitedAmount[1]
		}
	case uint32:
		left = strconv.FormatUint(uint64(amount/100), 10)
		rignt = strconv.FormatUint(uint64(amount%100), 10)
		if len(rignt) == 1 {
			rignt = "0" + rignt
		}
	case uint64:
		left = strconv.FormatUint(uint64(amount/100), 10)
		rignt = strconv.FormatUint(uint64(amount%100), 10)
		if len(rignt) == 1 {
			rignt = "0" + rignt
		}
	}

	return left, rignt
}

// OmmitEmptyStrings filters out empty strings from the provided slice of strings.
// It returns a new slice containing only the non-empty strings from the input.
//
// Parameters:
//
//   - strs: A slice of strings to be filtered.
//
// Returns:
//
//   - A new slice containing only the non-empty strings from the input.
func OmmitEmptyStrings(strs []string) []string {
	res := make([]string, 0, len(strs))
	for _, s := range strs {
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}

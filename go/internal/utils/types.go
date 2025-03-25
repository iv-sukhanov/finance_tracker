package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// UUIDsToStrings converts a slice of UUIDs to a slice of strings.
func UUIDsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, guid := range uuids {
		strs[i] = guid.String()
	}
	return strs
}

// ExtractAmountParts receives an amount in different formats (string, uint32, uint64)
// and returns the integer and fractional parts of the amount as strings.
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

// OmmitEmptyStrings removes empty strings from a slice of strings.
func OmmitEmptyStrings(strs []string) []string {
	res := make([]string, 0, len(strs))
	for _, s := range strs {
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}

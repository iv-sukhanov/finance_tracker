package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func UUIDsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, guid := range uuids {
		strs[i] = guid.String()
	}
	return strs
}

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

func OmmitEmptyStrings(strs []string) []string {
	res := make([]string, 0, len(strs))
	for _, s := range strs {
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}

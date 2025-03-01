package utils

import "github.com/google/uuid"

func UUIDsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, guid := range uuids {
		strs[i] = guid.String()
	}
	return strs
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

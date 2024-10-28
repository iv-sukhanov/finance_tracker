package typesh

import "github.com/google/uuid"

func UUIDsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, guid := range uuids {
		strs[i] = guid.String()
	}
	return strs
}

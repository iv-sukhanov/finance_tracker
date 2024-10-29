package utils

import (
	"fmt"
)

func MakeWhereIn(col string, fields ...string) string {
	if len(fields) == 0 {
		return ""
	}

	where := fmt.Sprintf(`WHERE (%s) IN ('`, col)
	for i, field := range fields {
		where += field
		if i != len(fields)-1 {
			where += "', '"
		} else {
			where += "')"
		}
	}
	return where
}

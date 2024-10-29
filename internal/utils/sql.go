package utils

import (
	"fmt"
	"time"
)

func MakeWhereIn(col string, op string, fields ...string) string {
	if len(fields) == 0 {
		return ""
	}

	where := fmt.Sprintf(`WHERE (%s) IN ('`, col)
	for i, field := range fields {
		where += field
		if i != len(fields)-1 {
			where += "', '"
		} else {
			where += "') " + op
		}
	}
	return where
}

func MakeIn(col string, op string, fields ...string) string {
	if len(fields) == 0 {
		return ""
	}

	where := fmt.Sprintf(`(%s) IN ('`, col)
	for i, field := range fields {
		where += field
		if i != len(fields)-1 {
			where += "', '"
		} else {
			where += "') " + op
		}
	}
	return where
}

func MakeLimit(limit int) string {
	if limit == -1 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", limit)
}

func MakeTimeFrame(col string, from, to time.Time, byTime bool) string {
	if !byTime {
		return ""
	}
	return fmt.Sprintf("WHERE %s BETWEEN %s AND %s", col, from, to)
}

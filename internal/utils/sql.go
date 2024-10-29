package utils

import (
	"fmt"
	"time"
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

func MakeIn(col string, op string, fields ...string) string {
	if len(fields) == 0 {
		return ""
	}

	where := fmt.Sprintf(`%s (%s) IN ('`, op, col)
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

func MakeLimit(limit int) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", limit)
}

func MakeTimeFrame(col, op string, from, to time.Time, byTime bool) string {
	if !byTime {
		return ""
	}
	return fmt.Sprintf("%s %s >= '%s' AND %s < '%s'", op, col, from.Format("2006-01-02 15:04:05"), col, to.Format("2006-01-02 15:04:05"))
}

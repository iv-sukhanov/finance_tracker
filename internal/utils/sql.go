package utils

import (
	"fmt"
	"strings"
	"time"
)

// delete
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

// could be rewritten with strings.Builder for consistency
func MakeIn(col string, fields ...string) string {
	if len(fields) == 0 {
		return ""
	}

	where := fmt.Sprintf(`(%s) IN ('`, col)
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

func MakeTimeFrame(col string, from, to time.Time, byTime bool) string {
	if !byTime {
		return ""
	}
	return fmt.Sprintf("%s >= '%s' AND %s < '%s'", col, from.Format("2006-01-02 15:04:05"), col, to.Format("2006-01-02 15:04:05"))
}

func BindWithOp(op string, needWhere bool, exprs ...string) string {
	if len(exprs) == 0 {
		return ""
	}
	nonEmptyExprs := OmmitEmptyStrings(exprs)

	var output strings.Builder
	if needWhere {
		output.WriteString("WHERE ")
	}

	for i, expr := range nonEmptyExprs {

		output.WriteString(expr)
		output.WriteRune(' ')
		if i != len(nonEmptyExprs)-1 {
			output.WriteString(op)
			output.WriteRune(' ')
		}
	}

	return output.String()
}

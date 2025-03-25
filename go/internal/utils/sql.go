package utils

import (
	"fmt"
	"strings"
	"time"
)

// MakeIn generates a SQL IN clause for the given column and fields.
// It returns an empty string if fields is empty.
//
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

// MakeLimit generates a SQL LIMIT clause for the given limit.
// It returns an empty string if limit is 0.
func MakeLimit(limit int) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", limit)
}

// MakeTimeFrame generates a condition for a column to be be between two timestamps.
// It returns an empty string if byTime is false.
func MakeTimeFrame(col string, from, to time.Time, byTime bool) string {
	if !byTime {
		return ""
	}
	return fmt.Sprintf("%s >= '%s' AND %s < '%s'", col, from.Format("2006-01-02 15:04:05"), col, to.Format("2006-01-02 15:04:05"))
}

// MakeOrderBy generates a SQL ORDER BY clause for the given column and order.
// It returns an empty string if col is empty.
func MakeOrderBy(col string, asc bool) string {

	if col == "" {
		return ""
	}

	var ascStr string
	if asc {
		ascStr = "ASC"
	} else {
		ascStr = "DESC"
	}
	return fmt.Sprintf("ORDER BY %s %s", col, ascStr)
}

// BindWithOp generates a SQL WHERE clause for the given expressions
// binding them with the given operator.
// It returns an empty string if exprs is empty or all expressions are empty.
func BindWithOp(op string, needWhere bool, exprs ...string) string {

	if len(exprs) == 0 {
		return ""
	}

	nonEmptyExprs := OmmitEmptyStrings(exprs)
	if len(nonEmptyExprs) == 0 {
		return ""
	}

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

// IsUniqueConstrainViolation checks if the error is a unique constraint violation.
func IsUniqueConstrainViolation(err error) bool {
	initialErr := GetItitialError(err)
	return GetSQLErrorCode(initialErr) == "23505"
}

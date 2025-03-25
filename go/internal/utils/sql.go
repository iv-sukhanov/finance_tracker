package utils

import (
	"fmt"
	"strings"
	"time"
)

// MakeIn constructs a SQL IN clause for a given column and a list of fields.
//
// Parameters:
//   - col: The name of the column to be used in the IN clause.
//   - fields: A variadic parameter containing the values to be included in the IN clause.
//
// Returns:
//
//	A string representing the constructed SQL clause. If no fields
//	are provided, an empty string is returned.
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
//
// Parameters:
//   - limit: The maximum number of rows to return.
//
// Returns:
//   - A string representing the SQL LIMIT clause. If limit is 0,
//     an empty string is returned.
func MakeLimit(limit int) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", limit)
}

// MakeTimeFrame generates a condition for a column to be be between two timestamps.
// It returns an empty string if byTime is false.
//
// Parameters:
//
//   - col: The name of the column to be used in the condition.
//   - from: The start time of the range.
//   - to: The end of the range
//   - byTime: A boolean indicating whether to include the time frame condition.
//
// Returns:
//   - A string representing the SQL condition. If byTime is false,
//     an empty string is returned.
func MakeTimeFrame(col string, from, to time.Time, byTime bool) string {
	if !byTime {
		return ""
	}
	return fmt.Sprintf("%s >= '%s' AND %s < '%s'", col, from.Format("2006-01-02 15:04:05"), col, to.Format("2006-01-02 15:04:05"))
}

// MakeOrderBy generates a SQL ORDER BY clause for the given column and order.
// It returns an empty string if col is empty.
//
// Parameters:
//   - col: The name of the column to be used in the ORDER BY clause.
//   - asc: A boolean indicating whether to sort in ascending order (true) or descending order (false).
//
// Returns:
//   - A string representing the SQL ORDER BY clause. If col is empty,
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

// BindWithOp constructs a SQL query string by combining multiple expressions
// with a specified logical operator (e.g., AND, OR). It optionally adds a
// "WHERE" clause at the beginning of the query.
//
// Parameters:
//   - op: The logical operator to use between expressions (e.g., "AND", "OR").
//   - needWhere: A boolean indicating whether to prepend "WHERE" to the query.
//   - exprs: A variadic parameter containing the expressions to combine.
//
// Returns:
//
//	A string representing the constructed SQL query. If no valid expressions
//	are provided, an empty string is returned.
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

// IsUniqueConstrainViolation checks if the provided error corresponds to a
// unique constraint violation in a SQL database. It extracts the initial
// error and compares its SQL error code to "23505", which is the standard
// code for unique constraint violations in PostgreSQL.
//
// Parameters:
//   - err: The error to be checked.
//
// Returns:
//   - bool: True if the error is a unique constraint violation, false otherwise.
func IsUniqueConstrainViolation(err error) bool {
	initialErr := GetItitialError(err)
	return GetSQLErrorCode(initialErr) == "23505"
}

package formatter

import (
	"fmt"
	"strings"
)

// FormatValue formats a value for SQL
func FormatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// FormatInClause formats an IN clause
func FormatInClause(tableName, field string, value interface{}) string {
	if values, ok := value.([]interface{}); ok {
		if len(values) == 0 {
			return fmt.Sprintf("FALSE /* empty IN clause for %s.%s */", tableName, field)
		}

		items := make([]string, len(values))
		for i, item := range values {
			items[i] = FormatValue(item)
		}
		return fmt.Sprintf("%s.%s IN (%s)", tableName, field, strings.Join(items, ", "))
	}
	return ""
}

// FormatEquality formats an equality condition
func FormatEquality(tableName, field string, value interface{}) string {
	return fmt.Sprintf("%s.%s = %s", tableName, field, FormatValue(value))
}

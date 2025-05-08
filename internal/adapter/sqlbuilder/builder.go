package sqlbuilder

import (
	"fmt"
	"strings"

	"mca-sql-convertor/internal/domain"
	"mca-sql-convertor/pkg/formatter"
)

// SQLBuilder converts domain queries to SQL
type SQLBuilder struct{}

// NewSQLBuilder creates a new SQLBuilder
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
}

// ConvertToSQL converts a domain query to SQL statements
func (b *SQLBuilder) ConvertToSQL(query *domain.Query) map[string]string {
	result := make(map[string]string)

	for tableName, tableQuery := range *query {
		// Build a combined query for the main table and its relations
		sql := b.buildCombinedSQL(tableName, tableQuery)
		result[tableName] = sql
	}

	return result
}

// buildCombinedSQL builds a single SQL query combining the main table and its relations
func (b *SQLBuilder) buildCombinedSQL(tableName string, query *domain.TableQuery) string {
	var sql strings.Builder

	// Get all related tables for joins
	joins := b.getJoinClauses(tableName, query)

	// Get all selected fields including from related tables
	selectedFields := b.getSelectedFields(tableName, query)

	// SELECT clause
	sql.WriteString("SELECT ")
	sql.WriteString(strings.Join(selectedFields, ", "))

	// FROM clause
	sql.WriteString(" FROM " + tableName)

	// JOIN clauses
	for _, join := range joins {
		sql.WriteString(" " + join)
	}

	// WHERE clause
	whereClause := b.buildWhereClause(tableName, query.Where)
	if whereClause != "" {
		sql.WriteString(" WHERE " + whereClause)
	}

	// ORDER BY clause
	orderClause := b.buildOrderClause(tableName, query.Order)
	if orderClause != "" {
		sql.WriteString(" ORDER BY " + orderClause)
	}

	// LIMIT clause
	if query.Limit != nil {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", *query.Limit))
	}

	return sql.String()
}

// getSelectedFields collects all selected fields from main table and relations
func (b *SQLBuilder) getSelectedFields(tableName string, query *domain.TableQuery) []string {
	// Start with fields from the main table
	var allFields []string

	if len(query.Select) > 0 {
		for _, field := range query.Select {
			allFields = append(allFields, tableName+"."+field)
		}
	} else {
		allFields = append(allFields, tableName+".*")
	}

	// Add fields from related tables
	for relationName, relationQuery := range query.Relations {
		if len(relationQuery.Select) > 0 {
			for _, field := range relationQuery.Select {
				allFields = append(allFields, relationName+"."+field)
			}
		}
	}

	return allFields
}

// getJoinClauses generates all JOIN clauses for related tables
func (b *SQLBuilder) getJoinClauses(tableName string, query *domain.TableQuery) []string {
	var joins []string

	// Process direct relations
	for relationName, relationQuery := range query.Relations {
		joinCondition := b.getJoinCondition(relationName, relationQuery, tableName)
		join := fmt.Sprintf("INNER JOIN %s ON %s", relationName, joinCondition)
		joins = append(joins, join)

		// Process nested relations recursively
		nestedJoins := b.getNestedJoinClauses(relationName, relationQuery)
		joins = append(joins, nestedJoins...)
	}

	return joins
}

// getNestedJoinClauses generates JOIN clauses for nested relations
func (b *SQLBuilder) getNestedJoinClauses(parentName string, parentQuery *domain.TableQuery) []string {
	var joins []string

	for relationName, relationQuery := range parentQuery.Relations {
		joinCondition := b.getJoinCondition(relationName, relationQuery, parentName)
		join := fmt.Sprintf("INNER JOIN %s ON %s", relationName, joinCondition)
		joins = append(joins, join)

		// Process further nested relations recursively
		nestedJoins := b.getNestedJoinClauses(relationName, relationQuery)
		joins = append(joins, nestedJoins...)
	}

	return joins
}

// getJoinCondition determines the join condition between tables
func (b *SQLBuilder) getJoinCondition(tableName string, query *domain.TableQuery, parentTable string) string {
	if query.Join != nil {
		parts := strings.Split(*query.Join, ":")
		if len(parts) == 2 {
			return fmt.Sprintf("%s.%s = %s.%s", tableName, parts[0], parentTable, parts[1])
		}
	}

	// Default join condition
	return fmt.Sprintf("%s.%s_id = %s.id", tableName, parentTable, parentTable)
}

// buildWhereClause builds the WHERE clause
func (b *SQLBuilder) buildWhereClause(tableName string, whereClause domain.WhereClause) string {
	var conditions []string

	// Process direct conditions
	for field, condition := range whereClause.Conditions {
		condSQL := b.buildCondition(tableName, field, condition)
		if condSQL != "" {
			conditions = append(conditions, condSQL)
		}
	}

	// Process AND conditions
	if len(whereClause.And) > 0 {
		var andConditions []string
		for _, andMap := range whereClause.And {
			for field, condition := range andMap {
				condSQL := b.buildCondition(tableName, field, condition)
				if condSQL != "" {
					andConditions = append(andConditions, condSQL)
				}
			}
		}
		if len(andConditions) > 0 {
			conditions = append(conditions, "("+strings.Join(andConditions, " AND ")+")")
		}
	}

	// Process OR conditions
	if len(whereClause.Or) > 0 {
		var orConditions []string
		for _, orMap := range whereClause.Or {
			for field, condition := range orMap {
				condSQL := b.buildCondition(tableName, field, condition)
				if condSQL != "" {
					orConditions = append(orConditions, condSQL)
				}
			}
		}
		if len(orConditions) > 0 {
			conditions = append(conditions, "("+strings.Join(orConditions, " OR ")+")")
		}
	}

	return strings.Join(conditions, " AND ")
}

// buildCondition builds a single condition
func (b *SQLBuilder) buildCondition(tableName, field string, condition interface{}) string {
	switch v := condition.(type) {
	case string, int, float64, bool:
		// Simple equality
		return formatter.FormatEquality(tableName, field, v)

	case map[string]interface{}:
		// Operator condition
		for op, value := range v {
			switch op {
			case ">":
				return fmt.Sprintf("%s.%s > %s", tableName, field, formatter.FormatValue(value))
			case ">=":
				return fmt.Sprintf("%s.%s >= %s", tableName, field, formatter.FormatValue(value))
			case "<":
				return fmt.Sprintf("%s.%s < %s", tableName, field, formatter.FormatValue(value))
			case "<=":
				return fmt.Sprintf("%s.%s <= %s", tableName, field, formatter.FormatValue(value))
			case "in":
				return formatter.FormatInClause(tableName, field, value)
				// Add other operators as needed
			}
		}
	}

	return ""
}

// buildOrderClause builds the ORDER BY clause
func (b *SQLBuilder) buildOrderClause(tableName string, orderValue interface{}) string {
	if orderValue == nil {
		return ""
	}

	// Handle string order
	if strOrder, ok := orderValue.(string); ok {
		direction := "ASC"
		field := strOrder
		if strings.HasPrefix(strOrder, "-") {
			direction = "DESC"
			field = strOrder[1:]
		}
		return fmt.Sprintf("%s.%s %s", tableName, field, direction)
	}

	// Handle array order
	if arrOrder, ok := orderValue.([]interface{}); ok {
		var orderClauses []string
		for _, item := range arrOrder {
			if strItem, isString := item.(string); isString {
				direction := "ASC"
				field := strItem
				if strings.HasPrefix(strItem, "-") {
					direction = "DESC"
					field = strItem[1:]
				}
				orderClauses = append(orderClauses, fmt.Sprintf("%s.%s %s", tableName, field, direction))
			}
		}
		return strings.Join(orderClauses, ", ")
	}

	return ""
}

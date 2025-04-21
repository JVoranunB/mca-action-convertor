package sqlbuilder

import (
	"fmt"
	"strings"

	"mca-action-convertor/internal/domain"
	"mca-action-convertor/pkg/formatter"
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
		sql := b.buildSQL(tableName, tableQuery, nil, "")
		result[tableName] = sql

		b.processRelations(result, tableName, tableQuery)
	}

	return result
}

// processRelations processes relation queries recursively
func (b *SQLBuilder) processRelations(result map[string]string, parentName string, parentQuery *domain.TableQuery) {
	for relationName, relationQuery := range parentQuery.Relations {
		fullName := parentName + "_" + relationName
		sql := b.buildSQL(relationName, relationQuery, parentQuery, parentName)
		result[fullName] = sql

		// Process nested relations
		b.processRelations(result, fullName, relationQuery)
	}
}

// buildSQL builds an SQL statement for a single table query
func (b *SQLBuilder) buildSQL(tableName string, query *domain.TableQuery, parentQuery *domain.TableQuery, parentTable string) string {
	var sql strings.Builder

	// SELECT clause
	sql.WriteString("SELECT ")
	if len(query.Select) > 0 {
		fields := make([]string, len(query.Select))
		for i, field := range query.Select {
			fields[i] = tableName + "." + field
		}
		sql.WriteString(strings.Join(fields, ", "))
	} else {
		sql.WriteString(tableName + ".*")
	}

	// FROM clause
	sql.WriteString("\nFROM " + tableName)

	// JOIN clause
	if parentQuery != nil {
		joinCondition := b.getJoinCondition(tableName, query, parentTable)
		sql.WriteString(fmt.Sprintf("\nINNER JOIN %s ON %s", parentTable, joinCondition))
	}

	// WHERE clause
	whereClause := b.buildWhereClause(tableName, query.Where)
	if whereClause != "" {
		sql.WriteString("\nWHERE " + whereClause)
	}

	// ORDER BY clause
	orderClause := b.buildOrderClause(tableName, query.Order)
	if orderClause != "" {
		sql.WriteString("\nORDER BY " + orderClause)
	}

	// LIMIT clause
	if query.Limit != nil {
		sql.WriteString(fmt.Sprintf("\nLIMIT %d", *query.Limit))
	}

	return sql.String()
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

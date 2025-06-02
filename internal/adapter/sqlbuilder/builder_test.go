package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mca-bigQuery/internal/domain"
)

func TestConvertToSQL(t *testing.T) {
	// Create a test query with relations
	query := createTestQuery()

	// Create SQL builder
	builder := NewSQLBuilder()

	// Convert query to SQL
	sqlMap := builder.ConvertToSQL(query)

	// Check if users query exists
	usersSQL, ok := sqlMap["users"]
	require.True(t, ok, "Expected 'users' query in SQL map")

	// Check if SQL contains expected parts (basic validation)
	expectedParts := []string{
		"SELECT users.id, users.username, users.email, posts.id, posts.title",
		"FROM users",
		"INNER JOIN posts ON posts.user_id = users.id",
		"WHERE users.status = 'active'",
		"ORDER BY users.username ASC",
		"LIMIT 10",
	}

	for _, part := range expectedParts {
		assert.Contains(t, usersSQL, part, "Expected SQL to contain '%s'", part)
	}

	// Verify that only one query is generated (combined query)
	assert.Len(t, sqlMap, 1, "Expected exactly one combined SQL query")
}

func TestBuildCombinedSQLWithMultipleRelations(t *testing.T) {
	// Create a query with multiple relations
	query := domain.Query{
		"orders": &domain.TableQuery{
			Select: []string{"id", "order_date", "total_amount"},
			Where: domain.WhereClause{
				Conditions: map[string]interface{}{
					"status": "completed",
				},
			},
			Relations: map[string]*domain.TableQuery{
				"customers": {
					Select: []string{"id", "name", "email"},
					Join:   domain.StrPtr("customer_id:id"),
				},
				"items": {
					Select: []string{"id", "product_id", "quantity"},
					Join:   domain.StrPtr("order_id:id"),
					Relations: map[string]*domain.TableQuery{
						"products": {
							Select: []string{"id", "name", "sku"},
							Join:   domain.StrPtr("id:product_id"),
						},
					},
				},
			},
		},
	}

	builder := NewSQLBuilder()
	sqlMap := builder.ConvertToSQL(&query)

	// Check if orders query exists
	ordersSQL, ok := sqlMap["orders"]
	require.True(t, ok, "Expected 'orders' query in SQL map")

	// Check if SQL contains all expected parts
	expectedParts := []string{
		"SELECT orders.id, orders.order_date, orders.total_amount, customers.id, customers.name, customers.email, items.id, items.product_id, items.quantity, products.id, products.name, products.sku",
		"FROM orders",
		"INNER JOIN customers ON customers.customer_id = orders.id",
		"INNER JOIN items ON items.order_id = orders.id",
		"INNER JOIN products ON products.id = items.product_id",
		"WHERE orders.status = 'completed'",
	}

	for _, part := range expectedParts {
		assert.Contains(t, ordersSQL, part, "Expected SQL to contain '%s'", part)
	}

	// Verify that only one query is generated (combined query)
	assert.Len(t, sqlMap, 1, "Expected exactly one combined SQL query")
}

func TestBuildWhereClauseWithComplexConditions(t *testing.T) {
	// Create a query with complex where clause
	query := domain.Query{
		"users": &domain.TableQuery{
			Select: []string{"id", "name"},
			Where: domain.WhereClause{
				And: []map[string]interface{}{
					{"status": "active"},
					{"created_at": map[string]interface{}{">": "2023-01-01"}},
				},
				Or: []map[string]interface{}{
					{"age": map[string]interface{}{">=": 18}},
					{"role": map[string]interface{}{"in": []interface{}{"admin", "editor"}}},
				},
			},
		},
	}

	builder := NewSQLBuilder()
	sqlMap := builder.ConvertToSQL(&query)

	usersSQL := sqlMap["users"]

	// Check for AND conditions
	assert.Contains(t, usersSQL, "(users.status = 'active' AND users.created_at > '2023-01-01')",
		"Expected SQL to contain AND conditions")

	// Check for OR conditions
	assert.Contains(t, usersSQL, "(users.age >= 18 OR users.role IN ('admin', 'editor'))",
		"Expected SQL to contain OR conditions")
}

func TestBuildOrderClause(t *testing.T) {
	testCases := []struct {
		name       string
		orderValue interface{}
		expected   string
	}{
		{
			name:       "Single field ascending",
			orderValue: "name",
			expected:   "table.name ASC",
		},
		{
			name:       "Single field descending",
			orderValue: "-created_at",
			expected:   "table.created_at DESC",
		},
		{
			name:       "Multiple fields",
			orderValue: []interface{}{"name", "-age"},
			expected:   "table.name ASC, table.age DESC",
		},
		{
			name:       "Empty order",
			orderValue: []interface{}{},
			expected:   "",
		},
		{
			name:       "Nil order",
			orderValue: nil,
			expected:   "",
		},
	}

	builder := NewSQLBuilder()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := builder.buildOrderClause("table", tc.orderValue)
			assert.Equal(t, tc.expected, result, "Order clause didn't match expected value")
		})
	}
}

func TestJoinConditionGeneration(t *testing.T) {
	builder := NewSQLBuilder()

	testCases := []struct {
		name        string
		tableName   string
		joinValue   *string
		parentTable string
		expected    string
	}{
		{
			name:        "Custom join",
			tableName:   "posts",
			joinValue:   domain.StrPtr("post_id:user_id"),
			parentTable: "users",
			expected:    "posts.post_id = users.user_id",
		},
		{
			name:        "Default join",
			tableName:   "comments",
			joinValue:   nil,
			parentTable: "posts",
			expected:    "comments.posts_id = posts.id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := &domain.TableQuery{
				Join: tc.joinValue,
			}

			joinCondition := builder.getJoinCondition(tc.tableName, query, tc.parentTable)
			assert.Equal(t, tc.expected, joinCondition, "Join condition didn't match expected value")
		})
	}
}

// Helper function to create a test query
func createTestQuery() *domain.Query {
	query := domain.Query{
		"users": &domain.TableQuery{
			Select: []string{"id", "username", "email"},
			Where: domain.WhereClause{
				Conditions: map[string]interface{}{
					"status": "active",
				},
			},
			Order: "username",
			Limit: domain.IntPtr(10),
			Relations: map[string]*domain.TableQuery{
				"posts": {
					Select: []string{"id", "title"},
					Where: domain.WhereClause{
						Conditions: map[string]interface{}{
							"published": true,
						},
					},
					Join: domain.StrPtr("user_id:id"),
				},
			},
		},
	}

	return &query
}

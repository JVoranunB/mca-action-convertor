package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mca-action-convertor/internal/domain"
)

func TestConvertToSQL(t *testing.T) {
	// Create a simple query
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
		"SELECT users.id, users.username, users.email",
		"FROM users",
		"WHERE users.status = 'active'",
		"ORDER BY users.username ASC",
		"LIMIT 10",
	}

	for _, part := range expectedParts {
		assert.Contains(t, usersSQL, part, "Expected SQL to contain '%s'", part)
	}

	// Check relation queries
	postsSQL, ok := sqlMap["users_posts"]
	require.True(t, ok, "Expected 'users_posts' relation query in SQL map")

	expectedPostsParts := []string{
		"SELECT posts.id, posts.title",
		"FROM posts",
		"INNER JOIN users ON posts.user_id = users.id",
	}

	for _, part := range expectedPostsParts {
		assert.Contains(t, postsSQL, part, "Expected relation SQL to contain '%s'", part)
	}
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

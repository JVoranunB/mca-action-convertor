package jsonparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJSON(t *testing.T) {
	// Setup
	parser := NewParser()

	// Test cases
	testCases := []struct {
		name    string
		jsonStr string
		wantErr bool
	}{
		{
			name: "Valid simple query",
			jsonStr: `{
				"users": {
					"select": ["id", "name"],
					"where": {
						"status": "active"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid complex query",
			jsonStr: `{
				"users": {
					"select": ["id", "username"],
					"where": {
						"and": [
							{ "status": "active" },
							{ "created_at": {">=": "2023-01-01"} }
						],
						"or": [
							{ "age": {">=": 18} },
							{ "role": {"in": ["admin", "editor"]} }
						]
					},
					"order": ["username"],
					"limit": 10
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid query with relations",
			jsonStr: `{
				"users": {
					"select": ["id", "name"],
					"posts": {
						"select": ["id", "title"]
					}
				}
			}`,
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			jsonStr: `{ this is not valid JSON }`,
			wantErr: true,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := parser.ParseJSON(tc.jsonStr)

			// Check error expectation
			if tc.wantErr {
				assert.Error(t, err, "Expected an error but didn't get one")
			} else {
				assert.NoError(t, err, "Expected no error but got one")

				// Verify query is not nil
				assert.NotNil(t, query, "ParseJSON() returned nil query")

				// Verify query has correct table
				assert.NotEmpty(t, *query, "ParseJSON() returned empty query")
			}
		})
	}
}

func TestWhereClauseUnmarshal(t *testing.T) {
	// Setup
	parser := NewParser()

	// Test case for where clause with direct conditions
	jsonStr := `{
		"users": {
			"where": {
				"status": "active",
				"age": 25
			}
		}
	}`

	// Parse JSON
	query, err := parser.ParseJSON(jsonStr)
	require.NoError(t, err, "Failed to parse JSON")

	// Verify direct conditions
	userQuery := (*query)["users"]
	require.NotNil(t, userQuery, "Failed to find 'users' in query")

	// Check conditions count
	assert.Len(t, userQuery.Where.Conditions, 2, "Expected 2 direct conditions")

	// Check specific condition values
	status, ok := userQuery.Where.Conditions["status"]
	assert.True(t, ok, "Expected 'status' condition to exist")
	assert.Equal(t, "active", status, "Expected 'status' condition to be 'active'")

	age, ok := userQuery.Where.Conditions["age"]
	assert.True(t, ok, "Expected 'age' condition to exist")
	assert.Equal(t, float64(25), age, "Expected 'age' condition to be 25")
}

func TestRelationsUnmarshal(t *testing.T) {
	// Setup
	parser := NewParser()

	// Test case for query with relations
	jsonStr := `{
		"users": {
			"select": ["id", "name"],
			"posts": {
				"select": ["id", "title"],
				"comments": {
					"select": ["id", "content"]
				}
			}
		}
	}`

	// Parse JSON
	query, err := parser.ParseJSON(jsonStr)
	require.NoError(t, err, "Failed to parse JSON")

	// Verify users table exists
	userQuery, ok := (*query)["users"]
	assert.True(t, ok, "Failed to find 'users' in query")

	// Verify posts relation exists
	assert.Contains(t, userQuery.Relations, "posts", "Failed to find 'posts' relation")

	// Get posts relation
	postsQuery := userQuery.Relations["posts"]
	assert.NotNil(t, postsQuery, "Posts relation is nil")

	// Verify posts select fields
	assert.ElementsMatch(t, []string{"id", "title"}, postsQuery.Select, "Posts select fields don't match")

	// Verify comments relation in posts
	assert.Contains(t, postsQuery.Relations, "comments", "Failed to find 'comments' relation")

	// Get comments relation
	commentsQuery := postsQuery.Relations["comments"]
	assert.NotNil(t, commentsQuery, "Comments relation is nil")

	// Verify comments select fields
	assert.ElementsMatch(t, []string{"id", "content"}, commentsQuery.Select, "Comments select fields don't match")
}

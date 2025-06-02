package test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mca-bigQuery/internal/adapter/jsonparser"
	"mca-bigQuery/internal/adapter/sqlbuilder"
	"mca-bigQuery/internal/repository"
	"mca-bigQuery/internal/usecase"
	"mca-bigQuery/test/setup"
)

// Test with sample file
func TestIntegrationWithFile(t *testing.T) {
	// Setup test environment
	env, err := setup.SetupStandardTestEnvironment()
	require.NoError(t, err, "Failed to setup test environment")
	defer env.Cleanup()

	// Initialize the components
	parser := jsonparser.NewParser()
	repo := repository.NewQueryRepository(parser)
	builder := sqlbuilder.NewSQLBuilder()
	converter := usecase.NewQueryConverterUseCase(repo, builder)

	// Setup our test files with known root table names
	// This ensures we map files to their expected root table names correctly
	fileToRootTable := map[string]string{
		"sample_query.json":  "users",  // This file has "users" as root
		"complex_query.json": "orders", // This file has "orders" as root
	}

	// Test both sample files
	for filename, expectedRootTable := range fileToRootTable {
		t.Run(filename, func(t *testing.T) {
			// Get full file path
			filePath := env.GetFullPath(filename)

			// Convert file to SQL
			sqlMap, err := converter.ConvertFileToSQL(filePath)
			require.NoError(t, err, "Failed to convert file %s", filename)
			require.NotEmpty(t, sqlMap, "No SQL statements generated from file %s", filename)

			// Check if the expected root table exists in results
			_, ok := sqlMap[expectedRootTable]
			assert.True(t, ok, "Expected root table '%s' not found in results for %s",
				expectedRootTable, filename)

			// Basic validation for each query
			for tableName, sql := range sqlMap {
				// Check that the SQL contains basic elements
				assert.Contains(t, sql, "SELECT", "SQL for %s doesn't contain SELECT", tableName)
				assert.Contains(t, sql, "FROM", "SQL for %s doesn't contain FROM", tableName)

				// If this is a relation table (contains underscore)
				if strings.Contains(tableName, "_") {
					assert.Contains(t, sql, "JOIN",
						"Relation SQL for %s doesn't contain JOIN", tableName)
				}

				t.Logf("Generated SQL for %s: %s", tableName, sql)
			}
		})
	}
}

// Function to help diagnose JSON file content
func TestPrintFileContent(t *testing.T) {
	// Only run this when debugging is needed
	if os.Getenv("DEBUG") != "1" {
		t.Skip("Skipping debug test")
	}

	// Setup test environment
	env, err := setup.SetupStandardTestEnvironment()
	require.NoError(t, err, "Failed to setup test environment")
	defer env.Cleanup()

	// Print content of each test file
	testFiles := []string{"sample_query.json", "complex_query.json"}
	for _, filename := range testFiles {
		filePath := env.GetFullPath(filename)
		content, err := os.ReadFile(filePath)
		require.NoError(t, err, "Failed to read file %s", filename)

		t.Logf("Content of %s:\n%s", filename, string(content))
	}
}

// The rest of your integration tests...

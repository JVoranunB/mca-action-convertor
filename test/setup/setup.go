package setup

import (
	"os"
	"path/filepath"
)

// TestEnvironment helps set up and clean up a test environment
type TestEnvironment struct {
	BaseDir   string
	TestFiles map[string]string
}

// NewTestEnvironment creates a new test environment
func NewTestEnvironment(baseDir string) *TestEnvironment {
	return &TestEnvironment{
		BaseDir:   baseDir,
		TestFiles: make(map[string]string),
	}
}

// AddTestFile adds a test file with the given content
func (e *TestEnvironment) AddTestFile(relativePath, content string) {
	e.TestFiles[relativePath] = content
}

// Setup creates the test environment
func (e *TestEnvironment) Setup() error {
	// Ensure base directory exists
	if err := os.MkdirAll(e.BaseDir, 0755); err != nil {
		return err
	}

	// Create all test files
	for path, content := range e.TestFiles {
		fullPath := filepath.Join(e.BaseDir, path)

		// Create directory for the file if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Write file content
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// Cleanup removes the test environment
func (e *TestEnvironment) Cleanup() error {
	// Remove the entire base directory
	return os.RemoveAll(e.BaseDir)
}

// GetFullPath returns the full path to a test file
func (e *TestEnvironment) GetFullPath(relativePath string) string {
	return filepath.Join(e.BaseDir, relativePath)
}

// SetupTestQueries sets up a standard set of test queries
func SetupStandardTestEnvironment() (*TestEnvironment, error) {
	env := NewTestEnvironment("test/testdata")

	// Add basic sample query
	env.AddTestFile("sample_query.json", `{
		"users": {
			"select": ["id", "username", "email"],
			"where": {
				"status": "active"
			},
			"limit": 10
		}
	}`)

	// Add complex query with relations
	env.AddTestFile("complex_query.json", `{
		"orders": {
			"select": ["id", "order_date", "total"],
			"where": {
				"status": "completed"
			},
			"customer": {
				"select": ["id", "name", "email"],
				"join": "customer_id:id"
			},
			"items": {
				"select": ["id", "product_id", "quantity"],
				"join": "order_id:id"
			}
		}
	}`)

	// Setup the environment
	if err := env.Setup(); err != nil {
		return nil, err
	}

	return env, nil
}

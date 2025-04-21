package domain

// Query represents a root query object that maps table names to TableQuery objects
type Query map[string]*TableQuery

// TableQuery represents the query for a single table
type TableQuery struct {
	Select    []string
	Where     WhereClause
	Order     interface{} // Can be string or []string
	Limit     *int
	Join      *string
	Relations map[string]*TableQuery
}

// Helper function to create int and string pointers
func IntPtr(i int) *int       { return &i }
func StrPtr(s string) *string { return &s }

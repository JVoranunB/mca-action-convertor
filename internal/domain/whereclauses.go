package domain

// WhereClause structure to handle logical operators
type WhereClause struct {
	And        []map[string]interface{}
	Or         []map[string]interface{}
	Conditions map[string]interface{} // For direct field conditions
}

// WhereOperator represents the various operators that can be used in where clauses
type WhereOperator string

const (
	OpEqual        WhereOperator = "="
	OpGreater      WhereOperator = ">"
	OpGreaterEqual WhereOperator = ">="
	OpLess         WhereOperator = "<"
	OpLessEqual    WhereOperator = "<="
	OpIn           WhereOperator = "in"
	// Add other operators as needed
)

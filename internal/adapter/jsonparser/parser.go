package jsonparser

import (
	"encoding/json"

	"mca-bigQuery/internal/domain"
)

// QueryDTO is a data transfer object for JSON unmarshaling
type QueryDTO map[string]*TableQueryDTO

// TableQueryDTO represents the JSON structure of a table query
type TableQueryDTO struct {
	Select    []string                  `json:"select,omitempty"`
	Where     WhereClauseDTO            `json:"where,omitempty"`
	Order     interface{}               `json:"order,omitempty"`
	Limit     *int                      `json:"limit,omitempty"`
	Join      *string                   `json:"join,omitempty"`
	Relations map[string]*TableQueryDTO `json:"-"` // Handled in custom unmarshaler
}

// WhereClauseDTO represents the JSON structure of where clauses
type WhereClauseDTO struct {
	And        []map[string]interface{} `json:"and,omitempty"`
	Or         []map[string]interface{} `json:"or,omitempty"`
	Conditions map[string]interface{}   `json:"-"` // For direct field conditions
}

// Parser provides methods to parse JSON into domain objects
type Parser struct{}

// NewParser creates a new JSON parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseJSON parses a JSON string into a domain Query
func (p *Parser) ParseJSON(jsonStr string) (*domain.Query, error) {
	var queryDTO QueryDTO
	if err := json.Unmarshal([]byte(jsonStr), &queryDTO); err != nil {
		return nil, err
	}

	query := mapDTOToDomain(queryDTO)
	return query, nil
}

// mapDTOToDomain converts DTO objects to domain objects
func mapDTOToDomain(queryDTO QueryDTO) *domain.Query {
	query := make(domain.Query)

	for tableName, tableQueryDTO := range queryDTO {
		query[tableName] = mapTableQueryDTOToDomain(tableQueryDTO)
	}

	return &query
}

// mapTableQueryDTOToDomain converts TableQueryDTO to domain TableQuery
func mapTableQueryDTOToDomain(dto *TableQueryDTO) *domain.TableQuery {
	tableQuery := &domain.TableQuery{
		Select:    dto.Select,
		Where:     mapWhereClauseDTOToDomain(dto.Where),
		Order:     dto.Order,
		Limit:     dto.Limit,
		Join:      dto.Join,
		Relations: make(map[string]*domain.TableQuery),
	}

	for relationName, relationDTO := range dto.Relations {
		tableQuery.Relations[relationName] = mapTableQueryDTOToDomain(relationDTO)
	}

	return tableQuery
}

// mapWhereClauseDTOToDomain converts WhereClauseDTO to domain WhereClause
func mapWhereClauseDTOToDomain(dto WhereClauseDTO) domain.WhereClause {
	return domain.WhereClause{
		And:        dto.And,
		Or:         dto.Or,
		Conditions: dto.Conditions,
	}
}

// Custom UnmarshalJSON for WhereClauseDTO
func (w *WhereClauseDTO) UnmarshalJSON(data []byte) error {
	// Try to unmarshal logical operators
	type LogicalOps struct {
		And []map[string]interface{} `json:"and,omitempty"`
		Or  []map[string]interface{} `json:"or,omitempty"`
	}

	var logicalOps LogicalOps
	if err := json.Unmarshal(data, &logicalOps); err != nil {
		return err
	}

	w.And = logicalOps.And
	w.Or = logicalOps.Or

	// Unmarshal all fields to capture direct conditions
	var allFields map[string]interface{}
	if err := json.Unmarshal(data, &allFields); err != nil {
		return err
	}

	// Initialize conditions map
	w.Conditions = make(map[string]interface{})

	// Add all fields except "and" and "or" as direct conditions
	for key, value := range allFields {
		if key != "and" && key != "or" {
			w.Conditions[key] = value
		}
	}

	return nil
}

// Custom UnmarshalJSON to handle nested relations
func (t *TableQueryDTO) UnmarshalJSON(data []byte) error {
	// First unmarshal standard fields
	type StandardFields struct {
		Select []string       `json:"select,omitempty"`
		Where  WhereClauseDTO `json:"where,omitempty"`
		Order  interface{}    `json:"order,omitempty"`
		Limit  *int           `json:"limit,omitempty"`
		Join   *string        `json:"join,omitempty"`
	}

	var std StandardFields
	if err := json.Unmarshal(data, &std); err != nil {
		return err
	}

	// Copy standard fields to our TableQueryDTO
	t.Select = std.Select
	t.Where = std.Where
	t.Order = std.Order
	t.Limit = std.Limit
	t.Join = std.Join

	// Now extract relations
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Initialize relations map
	t.Relations = make(map[string]*TableQueryDTO)

	// Standard fields to skip
	standardFields := map[string]bool{
		"select": true, "where": true, "order": true, "limit": true, "join": true,
	}

	// Process relations
	for key, value := range rawMap {
		if !standardFields[key] {
			var relation TableQueryDTO
			if err := json.Unmarshal(value, &relation); err != nil {
				return err
			}
			t.Relations[key] = &relation
		}
	}

	return nil
}

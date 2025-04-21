package repository

import (
	"io/ioutil"

	"mca-action-convertor/internal/domain"
)

// QueryRepository defines the interface for loading queries
type QueryRepository interface {
	ParseQuery(jsonStr string) (*domain.Query, error)
	LoadQueryFromFile(filename string) (*domain.Query, error)
}

// QueryParser defines the interface for parsing JSON queries
type QueryParser interface {
	ParseJSON(jsonStr string) (*domain.Query, error)
}

// QueryRepositoryImpl implements QueryRepository
type QueryRepositoryImpl struct {
	parser QueryParser
}

// NewQueryRepository creates a new QueryRepository
func NewQueryRepository(parser QueryParser) *QueryRepositoryImpl {
	return &QueryRepositoryImpl{
		parser: parser,
	}
}

// ParseQuery parses a query from a JSON string
func (r *QueryRepositoryImpl) ParseQuery(jsonStr string) (*domain.Query, error) {
	return r.parser.ParseJSON(jsonStr)
}

// LoadQueryFromFile loads a query from a file
func (r *QueryRepositoryImpl) LoadQueryFromFile(filename string) (*domain.Query, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return r.ParseQuery(string(data))
}

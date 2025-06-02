package usecase

import (
	"mca-bigQuery/internal/domain"
	"mca-bigQuery/internal/repository"
)

// SQLBuilderPort defines the interface for SQL building
type SQLBuilderPort interface {
	ConvertToSQL(query *domain.Query) map[string]string
}

// QueryConverterUseCase defines use cases for query conversion
type QueryConverterUseCase struct {
	repository repository.QueryRepository
	sqlBuilder SQLBuilderPort
}

// NewQueryConverterUseCase creates a new query converter use case
func NewQueryConverterUseCase(repo repository.QueryRepository, builder SQLBuilderPort) *QueryConverterUseCase {
	return &QueryConverterUseCase{
		repository: repo,
		sqlBuilder: builder,
	}
}

// ConvertJSONToSQL converts a JSON query string to SQL
func (uc *QueryConverterUseCase) ConvertJSONToSQL(jsonStr string) (map[string]string, error) {
	query, err := uc.repository.ParseQuery(jsonStr)
	if err != nil {
		return nil, err
	}

	return uc.sqlBuilder.ConvertToSQL(query), nil
}

// ConvertFileToSQL converts a query from a file to SQL
func (uc *QueryConverterUseCase) ConvertFileToSQL(filename string) (map[string]string, error) {
	query, err := uc.repository.LoadQueryFromFile(filename)
	if err != nil {
		return nil, err
	}

	return uc.sqlBuilder.ConvertToSQL(query), nil
}

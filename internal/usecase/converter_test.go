package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"mca-bigQuery/internal/domain"
)

// Define mocks using testify/mock
type MockQueryRepository struct {
	mock.Mock
}

func (m *MockQueryRepository) ParseQuery(jsonStr string) (*domain.Query, error) {
	args := m.Called(jsonStr)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Query), args.Error(1)
}

func (m *MockQueryRepository) LoadQueryFromFile(filename string) (*domain.Query, error) {
	args := m.Called(filename)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Query), args.Error(1)
}

type MockSQLBuilder struct {
	mock.Mock
}

func (m *MockSQLBuilder) ConvertToSQL(query *domain.Query) map[string]string {
	args := m.Called(query)
	return args.Get(0).(map[string]string)
}

// Test suite for QueryConverterUseCase
type ConverterTestSuite struct {
	suite.Suite
	mockRepo    *MockQueryRepository
	mockBuilder *MockSQLBuilder
	useCase     *QueryConverterUseCase
}

// Setup runs before each test
func (s *ConverterTestSuite) SetupTest() {
	s.mockRepo = new(MockQueryRepository)
	s.mockBuilder = new(MockSQLBuilder)
	s.useCase = NewQueryConverterUseCase(s.mockRepo, s.mockBuilder)
}

// Test ConvertJSONToSQL - success case
func (s *ConverterTestSuite) TestConvertJSONToSQL_Success() {
	// Setup
	jsonStr := `{"users":{"select":["id"]}}`

	queryResult := &domain.Query{
		"users": &domain.TableQuery{
			Select: []string{"id"},
		},
	}

	sqlResult := map[string]string{
		"users": "SELECT users.id FROM users",
	}

	// Configure mocks
	s.mockRepo.On("ParseQuery", jsonStr).Return(queryResult, nil)
	s.mockBuilder.On("ConvertToSQL", queryResult).Return(sqlResult)

	// Execute
	result, err := s.useCase.ConvertJSONToSQL(jsonStr)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), sqlResult, result)
	s.mockRepo.AssertExpectations(s.T())
	s.mockBuilder.AssertExpectations(s.T())
}

// Test ConvertJSONToSQL - parse error
func (s *ConverterTestSuite) TestConvertJSONToSQL_ParseError() {
	// Setup
	jsonStr := `invalid json`
	expectedErr := errors.New("parse error")

	// Configure mock
	s.mockRepo.On("ParseQuery", jsonStr).Return(nil, expectedErr)

	// Execute
	result, err := s.useCase.ConvertJSONToSQL(jsonStr)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedErr, err)
	assert.Nil(s.T(), result)
	s.mockRepo.AssertExpectations(s.T())
	// Builder should not be called if there's a parse error
	s.mockBuilder.AssertNotCalled(s.T(), "ConvertToSQL")
}

// Test ConvertFileToSQL - success case
func (s *ConverterTestSuite) TestConvertFileToSQL_Success() {
	// Setup
	filename := "test.json"

	queryResult := &domain.Query{
		"users": &domain.TableQuery{
			Select: []string{"id"},
		},
	}

	sqlResult := map[string]string{
		"users": "SELECT users.id FROM users",
	}

	// Configure mocks
	s.mockRepo.On("LoadQueryFromFile", filename).Return(queryResult, nil)
	s.mockBuilder.On("ConvertToSQL", queryResult).Return(sqlResult)

	// Execute
	result, err := s.useCase.ConvertFileToSQL(filename)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), sqlResult, result)
	s.mockRepo.AssertExpectations(s.T())
	s.mockBuilder.AssertExpectations(s.T())
}

// Test ConvertFileToSQL - file error
func (s *ConverterTestSuite) TestConvertFileToSQL_FileError() {
	// Setup
	filename := "nonexistent.json"
	expectedErr := errors.New("file not found")

	// Configure mock
	s.mockRepo.On("LoadQueryFromFile", filename).Return(nil, expectedErr)

	// Execute
	result, err := s.useCase.ConvertFileToSQL(filename)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedErr, err)
	assert.Nil(s.T(), result)
	s.mockRepo.AssertExpectations(s.T())
	// Builder should not be called if there's a file error
	s.mockBuilder.AssertNotCalled(s.T(), "ConvertToSQL")
}

// Run the test suite
func TestConverterTestSuite(t *testing.T) {
	suite.Run(t, new(ConverterTestSuite))
}

// Table-driven test example using testify/assert
func TestConvertJSONToSQL_TableDriven(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		jsonStr     string
		setupMocks  func(repo *MockQueryRepository, builder *MockSQLBuilder)
		expectErr   bool
		expectedSQL map[string]string
	}{
		{
			name:    "Simple Query",
			jsonStr: `{"users":{"select":["id","name"]}}`,
			setupMocks: func(repo *MockQueryRepository, builder *MockSQLBuilder) {
				queryResult := &domain.Query{
					"users": &domain.TableQuery{
						Select: []string{"id", "name"},
					},
				}
				sqlResult := map[string]string{
					"users": "SELECT users.id, users.name FROM users",
				}
				repo.On("ParseQuery", mock.Anything).Return(queryResult, nil)
				builder.On("ConvertToSQL", queryResult).Return(sqlResult)
			},
			expectErr: false,
			expectedSQL: map[string]string{
				"users": "SELECT users.id, users.name FROM users",
			},
		},
		{
			name:    "Parse Error",
			jsonStr: `invalid json`,
			setupMocks: func(repo *MockQueryRepository, builder *MockSQLBuilder) {
				repo.On("ParseQuery", mock.Anything).Return(nil, errors.New("parse error"))
			},
			expectErr:   true,
			expectedSQL: nil,
		},
		{
			name:    "Complex Query",
			jsonStr: `{"users":{"select":["id"],"posts":{"select":["title"]}}}`,
			setupMocks: func(repo *MockQueryRepository, builder *MockSQLBuilder) {
				queryResult := &domain.Query{
					"users": &domain.TableQuery{
						Select: []string{"id"},
						Relations: map[string]*domain.TableQuery{
							"posts": {
								Select: []string{"title"},
							},
						},
					},
				}
				sqlResult := map[string]string{
					"users":       "SELECT users.id FROM users",
					"users_posts": "SELECT posts.title FROM posts JOIN users ON posts.user_id = users.id",
				}
				repo.On("ParseQuery", mock.Anything).Return(queryResult, nil)
				builder.On("ConvertToSQL", queryResult).Return(sqlResult)
			},
			expectErr: false,
			expectedSQL: map[string]string{
				"users":       "SELECT users.id FROM users",
				"users_posts": "SELECT posts.title FROM posts JOIN users ON posts.user_id = users.id",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockQueryRepository)
			mockBuilder := new(MockSQLBuilder)
			useCase := NewQueryConverterUseCase(mockRepo, mockBuilder)

			// Configure mocks
			tc.setupMocks(mockRepo, mockBuilder)

			// Execute
			result, err := useCase.ConvertJSONToSQL(tc.jsonStr)

			// Assert
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedSQL, result)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
			mockBuilder.AssertExpectations(t)
		})
	}
}

// Table-driven test for file conversion
func TestConvertFileToSQL_TableDriven(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		filename    string
		setupMocks  func(repo *MockQueryRepository, builder *MockSQLBuilder)
		expectErr   bool
		expectedSQL map[string]string
	}{
		{
			name:     "Successful file conversion",
			filename: "test.json",
			setupMocks: func(repo *MockQueryRepository, builder *MockSQLBuilder) {
				queryResult := &domain.Query{
					"users": &domain.TableQuery{
						Select: []string{"id"},
					},
				}
				sqlResult := map[string]string{
					"users": "SELECT users.id FROM users",
				}
				repo.On("LoadQueryFromFile", "test.json").Return(queryResult, nil)
				builder.On("ConvertToSQL", queryResult).Return(sqlResult)
			},
			expectErr: false,
			expectedSQL: map[string]string{
				"users": "SELECT users.id FROM users",
			},
		},
		{
			name:     "File load error",
			filename: "nonexistent.json",
			setupMocks: func(repo *MockQueryRepository, builder *MockSQLBuilder) {
				repo.On("LoadQueryFromFile", "nonexistent.json").Return(nil, errors.New("file not found"))
			},
			expectErr:   true,
			expectedSQL: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockQueryRepository)
			mockBuilder := new(MockSQLBuilder)
			useCase := NewQueryConverterUseCase(mockRepo, mockBuilder)

			// Configure mocks
			tc.setupMocks(mockRepo, mockBuilder)

			// Execute
			result, err := useCase.ConvertFileToSQL(tc.filename)

			// Assert
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedSQL, result)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
			mockBuilder.AssertExpectations(t)
		})
	}
}

// Basic usage example
package main

import (
	"fmt"
	"log"

	"mca-sql-convertor/internal/adapter/jsonparser"
	"mca-sql-convertor/internal/adapter/sqlbuilder"
	"mca-sql-convertor/internal/repository"
	"mca-sql-convertor/internal/usecase"
)

func main() {
	// Initialize the components
	parser := jsonparser.NewParser()
	repo := repository.NewQueryRepository(parser)
	sqlBuilder := sqlbuilder.NewSQLBuilder()
	converter := usecase.NewQueryConverterUseCase(repo, sqlBuilder)

	// Example 1: Simple query
	simpleJSON := `{
		"users": {
			"select": ["id", "name", "email"],
			"where": {
				"active": true
			},
			"limit": 20
		}
	}`

	sqlMap, err := converter.ConvertJSONToSQL(simpleJSON)
	if err != nil {
		log.Fatalf("Error converting simple query: %v", err)
	}

	fmt.Println("=== Simple Query ===")
	for name, sql := range sqlMap {
		fmt.Printf("-- %s\n%s;\n\n", name, sql)
	}

	// Example 2: Complex query with relations
	complexJSON := `{
		"orders": {
			"select": ["id", "order_number", "customer_id", "total"],
			"where": {
				"and": [
					{ "status": "completed" },
					{ "created_at": {">=": "2023-01-01"} }
				],
				"or": [
					{ "total": {">": 100} },
					{ "priority": {"in": ["high", "urgent"]} }
				]
			},
			"order": ["-created_at", "total"],
			"limit": 10,
			"items": {
				"select": ["id", "product_id", "quantity", "price"],
				"join": "order_id:id",
				"where": {
					"quantity": {">": 0}
				},
				"products": {
					"select": ["id", "name", "sku"],
					"join": "id:product_id"
				}
			}
		}
	}`

	sqlMap, err = converter.ConvertJSONToSQL(complexJSON)
	if err != nil {
		log.Fatalf("Error converting complex query: %v", err)
	}

	fmt.Println("=== Complex Query with Relations ===")
	for name, sql := range sqlMap {
		fmt.Printf("-- %s\n%s;\n\n", name, sql)
	}

	// Example 3: Loading from file
	/*
		fmt.Println("=== Loading from File ===")
		fmt.Println("// Uncomment and adjust path as needed")

		sqlMap, err = converter.ConvertFileToSQL("test/testdata/sample_query.json")
		if err != nil {
			log.Fatalf("Error converting from file: %v", err)
		}

		for name, sql := range sqlMap {
			fmt.Printf("-- %s\n%s;\n\n", name, sql)
		}
	*/
}

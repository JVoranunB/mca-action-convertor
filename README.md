# MCA Action Converter

A RESTful API service that converts JSON-based queries into SQL statements, designed using Clean Architecture principles.

## Features

- Convert JSON queries to SQL statements via REST API
- Generate combined SQL queries with JOINs in a single statement
- Support for direct JSON format
- JSON validation and error handling
- Clean Architecture design

## Architecture

This project follows Clean Architecture principles with these layers:

1. **Domain Layer** - Core business entities
2. **Use Case Layer** - Application-specific business rules
3. **Interface Adapters Layer** - Converts data between layers
4. **Repository Layer** - Abstracts data access
5. **Delivery Layer** - HTTP handlers and routes
6. **Infrastructure Layer** - Configuration and logging

## API Endpoints

- `GET /api/v1/health` - Health check
- `POST /api/v1/convert` - Convert JSON query to SQL

## Running the Service

### Running Locally

To run the service locally:

```bash
# Install dependencies
go mod download

# Run the service
go run cmd/api/main.go
```

## Example Usage

### Convert JSON to SQL

```bash
curl -X POST http://localhost:3000/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "users": {
      "select": ["id", "username", "email", "created_at"],
      "where": {
        "and": [
          { "status": "active" },
          { "created_at": { ">=": "2023-01-01" } }
        ],
        "or": [
          { "age": { ">=": 18 } },
          { "role": { "in": ["admin", "editor"] } }
        ]
      },
      "order": ["username", "created_at"],
      "limit": 10,
      "orders": {
        "select": ["id"],
        "join": "user_id:id"
      }
    }
  }'
```

The resulting output will be a combined SQL query:

```json
{
  "status": "success",
  "data": {
    "queries": {
      "users": "SELECT users.id, users.username, users.email, users.created_at, orders.id\nFROM users\nINNER JOIN orders ON orders.user_id = users.id\nWHERE (users.status = 'active' AND users.created_at >= '2023-01-01') AND (users.age >= 18 OR users.role IN ('admin', 'editor'))\nORDER BY users.username ASC, users.created_at ASC\nLIMIT 10"
    }
  }
}
```

## JSON Query Format

The service accepts JSON queries in the following format:

```json
{
  "table_name": {
    "select": ["field1", "field2"], 
    "where": {
      "field": "value",
      "and": [
        { "field1": "value1" },
        { "field2": { ">": 10 } }
      ],
      "or": [
        { "field3": { "in": ["value1", "value2"] } },
        { "field4": "value4" }
      ]
    },
    "order": ["field", "-field2"],
    "limit": 10,
    "relation_table": {
      "select": ["field1"],
      "join": "foreign_key:primary_key" 
    }
  }
}
```

### Key Features of the Query Format

1. **Main Table**: The root object key defines the main table in the query.

2. **Select Fields**: The `select` array specifies which fields to select from the table.

3. **Where Clauses**:
   - Direct conditions: `"field": "value"`
   - AND conditions: `"and": [ { "field1": "value1" }, ... ]`
   - OR conditions: `"or": [ { "field1": "value1" }, ... ]`
   - Operators: `"field": { ">": value }`, `"field": { "in": [value1, value2] }`

4. **Order By**: 
   - String: `"order": "field"` (ascending) or `"order": "-field"` (descending)
   - Array: `"order": ["field1", "-field2"]`

5. **Limit**: `"limit": 10`

6. **Relations**: Nested objects represent related tables to join with
   - `"join": "foreign_key:primary_key"` specifies the join condition
   - If omitted, a default join condition is used: `relation.main_table_id = main_table.id`

### Complex Query Example

Here's a more complex example that generates a combined SQL query with joins:

```json
{
  "orders": {
    "select": ["id", "order_number", "total_amount", "status"],
    "where": {
      "and": [
        { "status": {"in": ["shipped", "delivered"]} },
        { "order_date": {">=": "2023-06-01"} }
      ],
      "total_amount": {">": 50}
    },
    "order": ["-order_date", "total_amount"],
    "limit": 100,
    "customer": {
      "select": ["id", "name", "email"],
      "join": "customer_id:id"
    },
    "items": {
      "select": ["id", "product_id", "quantity"],
      "join": "order_id:id",
      "products": {
        "select": ["id", "name", "sku"],
        "join": "id:product_id"
      }
    }
  }
}
```

This will generate a single SQL query with the appropriate JOINs.

## Project Structure

```
mca-bigQuery/
├── api/                       # API documentation
├── cmd/
│   └── api/                   # Application entry point
│       └── main.go
├── internal/
│   ├── adapter/               # Interface adapters
│   │   ├── jsonparser/
│   │   └── sqlbuilder/
│   ├── domain/                # Core domain entities
│   ├── handler/               # HTTP handlers
│   │   └── handler.go
│   ├── infrastructure/        # Infrastructure concerns
│   │   └── config/
│   │       ├── env.go         # Environment configuration
│   │       └── logger.go      # Logging configuration
│   ├── repository/            # Data access interfaces
│   ├── routes/                # API routes
│   │   └── routes.go
│   └── usecase/               # Business logic
├── pkg/                       # Shared utilities
│   └── formatter/
├── test/                      # Tests
└── README.md                  # This file
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

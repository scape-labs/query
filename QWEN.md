# QWEN.md - Project Context for Query Builder

## Project Overview

This is a Go library that provides a fluent SQL query builder for constructing various types of SQL queries (SELECT, INSERT, UPDATE, DELETE) with support for relationships between tables using JOIN operations. The library supports different parameter placeholder styles (question marks `?` or dollar-numbered `$1, $2`) and provides a chainable API for building queries.

The project is licensed under the MIT License, allowing for free use, modification, and distribution with proper attribution.

## Key Features

- Fluent interface for building SQL queries
- Support for SELECT, INSERT, UPDATE, and DELETE operations
- Parameter binding to prevent SQL injection
- Automatic escaping of table names, column names, and identifiers to prevent SQL injection
- Support for different parameter placeholder styles (QuestionMark `?` or DollarNumber `$1, $2`)
- WHERE clause construction with AND/OR conditions
- ORDER BY, LIMIT, and OFFSET support
- JOIN operations for handling table relationships
- Table alias support
- Chainable API design

## Project Structure

```
.
├── CLAUDE.md          # Guidance for Claude Code when working with this repository
├── go.mod             # Go module definition
├── LICENSE            # MIT License
├── query.go           # Main implementation file (single file library)
├── README.md          # Basic project description
└── .claude/
    └── settings.local.json  # Claude-specific settings
```

## Technology Stack

- **Language**: Go (Golang)
- **Module**: `github.com/scape-labs/query`
- **Go Version**: 1.24.1

## Core Components

### QueryBuilder (query.go)
The main struct that builds SQL queries using a fluent interface. Contains:
- `queryType`: Type of query (SelectQuery, InsertQuery, UpdateQuery, DeleteQuery)
- `table`: Target table name
- `tableAlias`: Alias for the main table
- `columns`: Selected columns (defaults to "*")
- `whereClauses`: Array of WHERE conditions
- `joinClauses`: Array of JOIN operations
- `order`, `limit`, `offset`: Query modifiers
- `paramStyle`: Parameter placeholder style (QuestionMark or DollarNumber)
- Fields for INSERT and UPDATE operations (insertColumns/insertValues, updateColumns/updateValues)

### WhereClause (query.go)
Represents individual WHERE conditions with:
- `Column`, `Operator`, `Value`: The condition components
- `JoinType`: How to join with other conditions ("AND"/"OR")

### JoinClause (query.go)
Represents individual JOIN operations with:
- `Type`: Type of JOIN (INNER, LEFT, RIGHT, FULL)
- `Table`: Table to join
- `Alias`: Alias for the joined table
- `Condition`: JOIN condition

### Query (query.go)
Final output containing the built SQL string and parameters array.

## Main API Methods

- `NewQueryBuilder()` - Creates new builder instance with default settings
- `ParameterPlaceholder(ParameterStyle)` - Sets parameter style (QuestionMark or DollarNumber)
- `Table(string)` - Sets target table
- `As(string)` - Sets table alias
- `Select(...string)` - Sets columns to select (for SELECT queries)
- `Insert(map[string]interface{})` - Sets data for INSERT query
- `InsertColumns(...string)` - Sets columns for INSERT query
- `Values(...interface{})` - Sets values for INSERT query
- `Update(map[string]interface{})` - Sets data for UPDATE query
- `Set(string, interface{})` - Sets individual column/value for UPDATE query
- `Delete()` - Sets query type to DELETE
- `Where(column, operator, value)` - Adds WHERE condition with AND logic
- `OrWhere(column, operator, value)` - Adds WHERE condition with OR logic
- `OrderBy(string)` - Sets ORDER BY clause
- `Limit(int)` / `Offset(int)` - Sets pagination
- `Join(table, condition)` - Adds JOIN clause
- `LeftJoin(table, condition)` - Adds LEFT JOIN clause
- `RightJoin(table, condition)` - Adds RIGHT JOIN clause
- `InnerJoin(table, condition)` - Adds INNER JOIN clause
- `FullJoin(table, condition)` - Adds FULL JOIN clause
- `JoinAs(table, alias, condition)` - Adds JOIN clause with table alias
- `LeftJoinAs(table, alias, condition)` - Adds LEFT JOIN clause with table alias
- `RightJoinAs(table, alias, condition)` - Adds RIGHT JOIN clause with table alias
- `InnerJoinAs(table, alias, condition)` - Adds INNER JOIN clause with table alias
- `FullJoinAs(table, alias, condition)` - Adds FULL JOIN clause with table alias
- `Build()` - Generates final Query with SQL and parameters

## Development Commands

### Go Commands
- `go mod tidy` - Download dependencies and clean up go.mod
- `go build .` - Build the library
- `go test .` - Run tests (if any exist)
- `go test -v .` - Run tests with verbose output
- `go fmt .` - Format Go code
- `go vet .` - Run Go vet for static analysis

### Module Management
- `go mod init github.com/scape-labs/query` - Initialize Go module (already done)
- `go get <package>` - Add new dependencies

## Key Design Patterns

1. **Fluent Interface**: All builder methods return `*QueryBuilder` for method chaining
2. **Parameter Binding**: Values are parameterized to prevent SQL injection
3. **Flexible Parameter Styles**: Supports both `?` and `$1, $2` placeholder formats
4. **Method Chaining**: Operations can be chained together for readable query construction
5. **Relationship Support**: JOIN operations enable querying related data across multiple tables
6. **Table Aliasing**: Support for table aliases improves query readability
7. **Security by Default**: Automatic escaping of identifiers prevents SQL injection

## Security Features

The query builder includes several built-in security features to prevent SQL injection:

1. **Parameter Binding**: All user data values are properly parameterized to prevent injection
2. **Identifier Escaping**: Table names, column names, and other SQL identifiers are automatically escaped
3. **Input Validation**: Dangerous SQL keywords and characters are filtered from identifiers
4. **Context-Aware Escaping**: Different escaping rules are applied based on the context (table names vs. WHERE clauses)

## Example Usage

```go
// SELECT query
qb := NewQueryBuilder().Table("users").Select("id", "name", "email").Where("age", ">", 18)
query := qb.Build()
// Result: SELECT id, name, email FROM users WHERE age > $1

// INSERT query
data := map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
}
qb := NewQueryBuilder().Table("users").Insert(data)
query := qb.Build()
// Result: INSERT INTO users (name, email, age) VALUES ($1, $2, $3)

// UPDATE query
qb := NewQueryBuilder().Table("users").Update(map[string]interface{}{
    "name": "Jane Doe",
    "email": "jane@example.com",
}).Where("id", "=", 1)
query := qb.Build()
// Result: UPDATE users SET name = $1, email = $2 WHERE id = $3

// DELETE query
qb := NewQueryBuilder().Table("users").Delete().Where("id", "=", 1)
query := qb.Build()
// Result: DELETE FROM users WHERE id = $1

// JOIN query with relationships
qb := NewQueryBuilder().
    Table("users").
    Select("users.id", "users.name", "accounts.name as account_name").
    LeftJoin("accounts", "accounts.id = users.account_id").
    Where("users.active", "=", true)
query := qb.Build()
// Result: SELECT users.id, users.name, accounts.name as account_name 
//         FROM users LEFT JOIN accounts on accounts.id = users.account_id 
//         WHERE users.active = $1

// Multiple JOINs with table aliases
qb := NewQueryBuilder().
    Table("posts").
    As("p").
    Select("p.title", "u.name as author", "c.name as category").
    LeftJoinAs("users", "u", "u.id = p.user_id").
    InnerJoinAs("categories", "c", "c.id = p.category_id").
    Where("p.published", "=", true).
    OrderBy("p.created_at").
    Limit(10)
query := qb.Build()
// Result: SELECT p.title, u.name as author, c.name as category 
//         FROM posts as p 
//         LEFT JOIN users as u on u.id = p.user_id 
//         INNER JOIN categories as c on c.id = p.category_id 
//         WHERE p.published = $1 
//         ORDER BY p.created_at 
//         LIMIT 10
```

## Notes and Considerations

1. The package declaration in `query.go` is `package query` while the module is `github.com/scape-labs/query`
2. This is a single-file library implementation
3. No test files currently exist in the repository
4. The library is designed to be lightweight and focused on SQL query building
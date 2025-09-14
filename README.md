# query
Simple Fluent Query Builder

A Go library for building SQL queries with a fluent interface. Supports SELECT, INSERT, UPDATE, and DELETE operations, as well as JOIN operations for handling relationships between tables.

## Features

- Fluent interface for building SQL queries
- Support for SELECT, INSERT, UPDATE, and DELETE operations
- Parameter binding to prevent SQL injection
- Support for different parameter placeholder styles (QuestionMark `?` or DollarNumber `$1, $2`)
- WHERE clause construction with AND/OR conditions
- ORDER BY, LIMIT, and OFFSET support
- JOIN operations for handling table relationships
- Table alias support
- Chainable API design

## Installation

```bash
go get github.com/scape-labs/query
```

## Usage

### Basic Query Building

```go
import "github.com/scape-labs/query"

// SELECT query
qb := query.NewQueryBuilder().
    Table("users").
    Select("id", "name", "email").
    Where("age", ">", 18).
    OrderBy("name").
    Limit(10)

query := qb.Build()
// Result: SELECT id, name, email FROM users WHERE age > $1 ORDER BY name LIMIT 10
```

### INSERT Operation

```go
data := map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
}

qb := query.NewQueryBuilder().
    Table("users").
    Insert(data)

query := qb.Build()
// Result: INSERT INTO users (name, email, age) VALUES ($1, $2, $3)
```

### UPDATE Operation

```go
qb := query.NewQueryBuilder().
    Table("users").
    Update(map[string]interface{}{
        "name": "Jane Doe",
        "email": "jane@example.com",
    }).
    Where("id", "=", 1)

query := qb.Build()
// Result: UPDATE users SET name = $1, email = $2 WHERE id = $3
```

### DELETE Operation

```go
qb := query.NewQueryBuilder().
    Table("users").
    Delete().
    Where("id", "=", 1)

query := qb.Build()
// Result: DELETE FROM users WHERE id = $1
```

### JOIN Operations (Relationships)

The query builder supports various JOIN operations to handle relationships between tables:

```go
// Simple LEFT JOIN
qb := query.NewQueryBuilder().
    Table("users").
    Select("users.id", "users.name", "accounts.name as account_name").
    LeftJoin("accounts", "accounts.id = users.account_id").
    Where("users.active", "=", true)

query := qb.Build()
// Result: SELECT users.id, users.name, accounts.name as account_name 
//         FROM users LEFT JOIN accounts on accounts.id = users.account_id 
//         WHERE users.active = $1
```

```go
// Multiple JOINs with table aliases
qb := query.NewQueryBuilder().
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

### Parameter Styles

The query builder supports different parameter placeholder styles:

```go
// Default is DollarNumber ($1, $2, etc.)
qb := query.NewQueryBuilder().
    Table("users").
    Where("id", "=", 1)

query := qb.Build()
// Result: SELECT * FROM users WHERE id = $1

// Switch to QuestionMark (?)
qb = query.NewQueryBuilder().
    ParameterPlaceholder(query.QuestionMark).
    Table("users").
    Where("id", "=", 1)

query = qb.Build()
// Result: SELECT * FROM users WHERE id = ?
```

## API Reference

### QueryBuilder Methods

- `NewQueryBuilder()` - Creates a new query builder instance
- `Table(name string)` - Sets the table name
- `As(alias string)` - Sets a table alias
- `Select(columns ...string)` - Sets the columns to select
- `Insert(data map[string]interface{})` - Sets data for INSERT operation
- `Update(data map[string]interface{})` - Sets data for UPDATE operation
- `Delete()` - Sets query type to DELETE
- `Where(column, operator string, value interface{})` - Adds a WHERE condition
- `OrWhere(column, operator string, value interface{})` - Adds an OR WHERE condition
- `OrderBy(order string)` - Sets the ORDER BY clause
- `Limit(limit int)` - Sets the LIMIT clause
- `Offset(offset int)` - Sets the OFFSET clause
- `ParameterPlaceholder(style ParameterStyle)` - Sets the parameter placeholder style

### JOIN Methods

- `Join(table, condition string)` - Adds a JOIN clause
- `LeftJoin(table, condition string)` - Adds a LEFT JOIN clause
- `RightJoin(table, condition string)` - Adds a RIGHT JOIN clause
- `InnerJoin(table, condition string)` - Adds an INNER JOIN clause
- `FullJoin(table, condition string)` - Adds a FULL JOIN clause
- `JoinAs(table, alias, condition string)` - Adds a JOIN clause with table alias
- `LeftJoinAs(table, alias, condition string)` - Adds a LEFT JOIN clause with table alias
- `RightJoinAs(table, alias, condition string)` - Adds a RIGHT JOIN clause with table alias
- `InnerJoinAs(table, alias, condition string)` - Adds an INNER JOIN clause with table alias
- `FullJoinAs(table, alias, condition string)` - Adds a FULL JOIN clause with table alias

### Types

- `ParameterStyle` - Enum for parameter placeholder styles (QuestionMark, DollarNumber)
- `QueryType` - Enum for query types (SelectQuery, InsertQuery, UpdateQuery, DeleteQuery)

### Structs

- `Query` - Contains the built SQL string and parameters
- `QueryBuilder` - The main struct for building queries
- `WhereClause` - Represents a WHERE condition
- `JoinClause` - Represents a JOIN operation
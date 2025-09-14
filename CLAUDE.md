# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library that provides a fluent SQL query builder for constructing various types of SQL queries (SELECT, INSERT, UPDATE, DELETE) with support for relationships between tables using JOIN operations. The library supports different parameter placeholder styles (question marks `?` or dollar-numbered `$1, $2`) and provides a chainable API for building queries.

## Commands

### Development Commands
- `go mod tidy` - Download dependencies and clean up go.mod
- `go build .` - Build the library
- `go run .` - Run the main package (if applicable)
- `go test .` - Run tests
- `go test -v .` - Run tests with verbose output
- `go fmt .` - Format Go code
- `go vet .` - Run Go vet for static analysis

### Module Management
- `go mod init github.com/scape-labs/query` - Initialize Go module (already done)
- `go get <package>` - Add new dependencies

## Architecture

### Core Components

**QueryBuilder (query.go)**: The main struct that builds SQL queries using a fluent interface. Contains:
- `queryType`: Type of query (SelectQuery, InsertQuery, UpdateQuery, DeleteQuery)
- `table`: Target table name
- `tableAlias`: Alias for the main table
- `columns`: Selected columns (defaults to "*")
- `whereClauses`: Array of WHERE conditions
- `joinClauses`: Array of JOIN operations
- `order`, `limit`, `offset`: Query modifiers
- `paramStyle`: Parameter placeholder style (QuestionMark or DollarNumber)
- Fields for INSERT and UPDATE operations (insertColumns/insertValues, updateColumns/updateValues)

**WhereClause (query.go)**: Represents individual WHERE conditions with:
- `Column`, `Operator`, `Value`: The condition components
- `JoinType`: How to join with other conditions (\"AND\"/\"OR\")

**JoinClause (query.go)**: Represents individual JOIN operations with:
- `Type`: Type of JOIN (INNER, LEFT, RIGHT, FULL)
- `Table`: Table to join
- `Alias`: Alias for the joined table
- `Condition`: JOIN condition

**Query (query.go)**: Final output containing the built SQL string and parameters array.

### Key Design Patterns

1. **Fluent Interface**: All builder methods return `*QueryBuilder` for method chaining
2. **Parameter Binding**: Values are parameterized to prevent SQL injection
3. **Flexible Parameter Styles**: Supports both `?` and `$1, $2` placeholder formats
4. **Relationship Support**: JOIN operations enable querying related data across multiple tables
5. **Table Aliasing**: Support for table aliases improves query readability

### Main API Methods
- `NewQueryBuilder()` - Creates new builder instance
- `Table(string)` - Sets target table
- `As(string)` - Sets table alias
- `Select(...string)` - Sets columns to select
- `Where(column, operator, value)` - Adds WHERE condition
- `OrWhere(column, operator, value)` - Adds OR WHERE condition
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
- `ParameterPlaceholder(ParameterStyle)` - Sets parameter style
- `Build()` - Generates final Query with SQL and parameters

## File Structure

- `query.go` - Main implementation (single file library)
- `go.mod` - Go module definition
- `README.md` - Basic project description
- `LICENSE` - Project license

## Package Declaration

The code uses `package query` which matches the module `github.com/scape-labs/query`.
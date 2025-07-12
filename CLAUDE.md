# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library that provides a fluent SQL query builder for constructing SELECT queries. The library supports different parameter placeholder styles (question marks `?` or dollar-numbered `$1, $2`) and provides a chainable API for building queries.

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

**QueryBuilder (query.go:20-28)**: The main struct that builds SQL queries using a fluent interface. Contains:
- `table`: Target table name
- `columns`: Selected columns (defaults to "*")
- `whereClauses`: Array of WHERE conditions
- `order`, `limit`, `offset`: Query modifiers
- `paramStyle`: Parameter placeholder style (QuestionMark or DollarNumber)

**WhereClause (query.go:30-35)**: Represents individual WHERE conditions with:
- `Column`, `Operator`, `Value`: The condition components
- `JoinType`: How to join with other conditions ("AND"/"OR")

**Query (query.go:15-18)**: Final output containing the built SQL string and parameters array.

### Key Design Patterns

1. **Fluent Interface**: All builder methods return `*QueryBuilder` for method chaining
2. **Parameter Binding**: Values are parameterized to prevent SQL injection
3. **Flexible Parameter Styles**: Supports both `?` and `$1, $2` placeholder formats

### Main API Methods
- `NewQueryBuilder()` - Creates new builder instance
- `Table(string)` - Sets target table
- `Select(...string)` - Sets columns to select
- `Where(column, operator, value)` - Adds WHERE condition
- `OrWhere(column, operator, value)` - Adds OR WHERE condition
- `OrderBy(string)` - Sets ORDER BY clause
- `Limit(int)` / `Offset(int)` - Sets pagination
- `ParameterPlaceholder(ParameterStyle)` - Sets parameter style
- `Build()` - Generates final Query with SQL and parameters

## File Structure

- `query.go` - Main implementation (single file library)
- `go.mod` - Go module definition
- `README.md` - Basic project description
- `LICENSE` - Project license

## Package Declaration

The code uses `package repositories` but the module is `github.com/scape-labs/query`. This suggests the library is intended to be used as part of a larger repositories package or the package declaration may need updating.
package query

import (
	"fmt"
	"strings"
)

type ParameterStyle int

const (
	QuestionMark ParameterStyle = iota // ?
	DollarNumber                       // $1, $2, etc.
)

type QueryType int

const (
	SelectQuery QueryType = iota
	InsertQuery
	UpdateQuery
	DeleteQuery
)

type Query struct {
	SQL    string
	Params []interface{}
}

func (q Query) Sql() string {
	return q.SQL
}

type QueryBuilder struct {
	queryType    QueryType
	table        string
	tableAlias   string
	columns      []string
	whereClauses []*WhereClause
	joinClauses  []*JoinClause
	order        string
	limit        int
	offset       int
	paramStyle   ParameterStyle

	// For INSERT operations
	insertColumns []string
	insertValues  []interface{}

	// For UPDATE operations
	updateColumns []string
	updateValues  []interface{}
}

type WhereClause struct {
	Column   string
	Operator string
	Value    interface{}
	JoinType string // AND/OR
}

// JoinClause represents a JOIN operation in a query
type JoinClause struct {
	Type      string // INNER, LEFT, RIGHT, FULL
	Table     string
	Alias     string
	Condition string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		queryType:   SelectQuery,
		columns:     []string{"*"},
		joinClauses: []*JoinClause{},
		paramStyle:  DollarNumber, // Default to DollarNumber
	}
}

func (b *QueryBuilder) ParameterPlaceholder(style ParameterStyle) *QueryBuilder {
	b.paramStyle = style
	return b
}

func (b *QueryBuilder) Table(table string) *QueryBuilder {
	b.table = table
	return b
}

// SELECT operations
func (b *QueryBuilder) Select(columns ...string) *QueryBuilder {
	b.queryType = SelectQuery
	if len(columns) > 0 {
		b.columns = columns
	}
	return b
}

// INSERT operations
func (b *QueryBuilder) Insert(data map[string]interface{}) *QueryBuilder {
	b.queryType = InsertQuery
	b.insertColumns = make([]string, 0, len(data))
	b.insertValues = make([]interface{}, 0, len(data))

	for column, value := range data {
		b.insertColumns = append(b.insertColumns, column)
		b.insertValues = append(b.insertValues, value)
	}
	return b
}

func (b *QueryBuilder) InsertColumns(columns ...string) *QueryBuilder {
	b.queryType = InsertQuery
	b.insertColumns = columns
	return b
}

func (b *QueryBuilder) Values(values ...interface{}) *QueryBuilder {
	b.insertValues = values
	return b
}

// UPDATE operations
func (b *QueryBuilder) Update(data map[string]interface{}) *QueryBuilder {
	b.queryType = UpdateQuery
	b.updateColumns = make([]string, 0, len(data))
	b.updateValues = make([]interface{}, 0, len(data))

	for column, value := range data {
		b.updateColumns = append(b.updateColumns, column)
		b.updateValues = append(b.updateValues, value)
	}
	return b
}

func (b *QueryBuilder) Set(column string, value interface{}) *QueryBuilder {
	b.queryType = UpdateQuery
	b.updateColumns = append(b.updateColumns, column)
	b.updateValues = append(b.updateValues, value)
	return b
}

// DELETE operations
func (b *QueryBuilder) Delete() *QueryBuilder {
	b.queryType = DeleteQuery
	return b
}

// WHERE clauses (common to all query types)
func (b *QueryBuilder) Where(column string, operator string, value interface{}) *QueryBuilder {
	b.whereClauses = append(b.whereClauses, &WhereClause{
		Column:   column,
		Operator: operator,
		Value:    value,
		JoinType: "and",
	})
	return b
}

func (b *QueryBuilder) OrWhere(column string, operator string, value interface{}) *QueryBuilder {
	b.whereClauses = append(b.whereClauses, &WhereClause{
		Column:   column,
		Operator: operator,
		Value:    value,
		JoinType: "or",
	})
	return b
}

// ORDER BY (for SELECT and UPDATE/DELETE with LIMIT support in some databases)
func (b *QueryBuilder) OrderBy(order string) *QueryBuilder {
	b.order = order
	return b
}

// LIMIT and OFFSET (primarily for SELECT, but some databases support for UPDATE/DELETE)
func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
	b.limit = limit
	return b
}

func (b *QueryBuilder) Offset(offset int) *QueryBuilder {
	b.offset = offset
	return b
}

// JOIN operations
func (b *QueryBuilder) Join(table, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "JOIN",
		Table:     table,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "LEFT JOIN",
		Table:     table,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) RightJoin(table, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "RIGHT JOIN",
		Table:     table,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) InnerJoin(table, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "INNER JOIN",
		Table:     table,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) FullJoin(table, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "FULL JOIN",
		Table:     table,
		Condition: condition,
	})
	return b
}

// JOIN operations with alias support
func (b *QueryBuilder) JoinAs(table, alias, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "JOIN",
		Table:     table,
		Alias:     alias,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) LeftJoinAs(table, alias, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "LEFT JOIN",
		Table:     table,
		Alias:     alias,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) RightJoinAs(table, alias, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "RIGHT JOIN",
		Table:     table,
		Alias:     alias,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) InnerJoinAs(table, alias, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "INNER JOIN",
		Table:     table,
		Alias:     alias,
		Condition: condition,
	})
	return b
}

func (b *QueryBuilder) FullJoinAs(table, alias, condition string) *QueryBuilder {
	b.joinClauses = append(b.joinClauses, &JoinClause{
		Type:      "FULL JOIN",
		Table:     table,
		Alias:     alias,
		Condition: condition,
	})
	return b
}

// Table alias support
func (b *QueryBuilder) As(alias string) *QueryBuilder {
	b.tableAlias = alias
	return b
}

func (b *QueryBuilder) getPlaceholder(index int) string {
	switch b.paramStyle {
	case QuestionMark:
		return "?"
	case DollarNumber:
		return fmt.Sprintf("$%d", index)
	default:
		return fmt.Sprintf("$%d", index) // Default to DollarNumber
	}
}

func (b *QueryBuilder) Build() Query {
	switch b.queryType {
	case SelectQuery:
		return b.buildSelect()
	case InsertQuery:
		return b.buildInsert()
	case UpdateQuery:
		return b.buildUpdate()
	case DeleteQuery:
		return b.buildDelete()
	default:
		return b.buildSelect()
	}
}

func (b *QueryBuilder) buildSelect() Query {
	var query strings.Builder
	var params []interface{}
	paramCount := 0

	// Build SELECT clause
	query.WriteString("select ")
	query.WriteString(strings.Join(b.columns, ", "))

	// Build FROM clause
	query.WriteString(" from ")
	query.WriteString(b.table)
	if b.tableAlias != "" {
		query.WriteString(" as ")
		query.WriteString(b.tableAlias)
	}

	// Build JOIN clauses
	for _, join := range b.joinClauses {
		query.WriteString(" ")
		query.WriteString(join.Type)
		query.WriteString(" ")
		query.WriteString(join.Table)
		if join.Alias != "" {
			query.WriteString(" as ")
			query.WriteString(join.Alias)
		}
		query.WriteString(" on ")
		query.WriteString(join.Condition)
	}

	// Build WHERE clause
	if len(b.whereClauses) > 0 {
		whereSQL, whereParams, count := b.buildWhereClause(paramCount)
		query.WriteString(whereSQL)
		params = append(params, whereParams...)
		paramCount = count
	}

	// Build ORDER BY clause
	if b.order != "" {
		query.WriteString(" order by ")
		query.WriteString(b.order)
	}

	// Build LIMIT clause
	if b.limit > 0 {
		query.WriteString(fmt.Sprintf(" limit %d", b.limit))
	}

	// Build OFFSET clause
	if b.offset > 0 {
		query.WriteString(fmt.Sprintf(" offset %d", b.offset))
	}

	return Query{
		SQL:    query.String(),
		Params: params,
	}
}

func (b *QueryBuilder) buildInsert() Query {
	var query strings.Builder
	var params []interface{}

	// Build INSERT clause
	query.WriteString("insert into ")
	query.WriteString(b.table)

	if len(b.insertColumns) > 0 {
		// Build columns
		query.WriteString(" (")
		query.WriteString(strings.Join(b.insertColumns, ", "))
		query.WriteString(") values (")

		// Build placeholders
		placeholders := make([]string, len(b.insertValues))
		for i := range b.insertValues {
			placeholders[i] = b.getPlaceholder(i + 1)
		}
		query.WriteString(strings.Join(placeholders, ", "))
		query.WriteString(")")

		params = append(params, b.insertValues...)
	}

	return Query{
		SQL:    query.String(),
		Params: params,
	}
}

func (b *QueryBuilder) buildUpdate() Query {
	var query strings.Builder
	var params []interface{}
	paramCount := 0

	// Build UPDATE clause
	query.WriteString("update ")
	query.WriteString(b.table)
	query.WriteString(" set ")

	// Build SET clause
	setClauses := make([]string, len(b.updateColumns))
	for i, column := range b.updateColumns {
		paramCount++
		setClauses[i] = fmt.Sprintf("%s = %s", column, b.getPlaceholder(paramCount))
	}
	query.WriteString(strings.Join(setClauses, ", "))
	params = append(params, b.updateValues...)

	// Build WHERE clause
	if len(b.whereClauses) > 0 {
		whereSQL, whereParams, count := b.buildWhereClause(paramCount)
		query.WriteString(whereSQL)
		params = append(params, whereParams...)
		paramCount = count
	}

	// Build ORDER BY clause (supported in some databases like MySQL)
	if b.order != "" {
		query.WriteString(" order by ")
		query.WriteString(b.order)
	}

	// Build LIMIT clause (supported in some databases like MySQL)
	if b.limit > 0 {
		query.WriteString(fmt.Sprintf(" limit %d", b.limit))
	}

	return Query{
		SQL:    query.String(),
		Params: params,
	}
}

func (b *QueryBuilder) buildDelete() Query {
	var query strings.Builder
	var params []interface{}
	paramCount := 0

	// Build DELETE clause
	query.WriteString("delete from ")
	query.WriteString(b.table)

	// Build WHERE clause
	if len(b.whereClauses) > 0 {
		whereSQL, whereParams, count := b.buildWhereClause(paramCount)
		query.WriteString(whereSQL)
		params = append(params, whereParams...)
		paramCount = count
	}

	// Build ORDER BY clause (supported in some databases like MySQL)
	if b.order != "" {
		query.WriteString(" order by ")
		query.WriteString(b.order)
	}

	// Build LIMIT clause (supported in some databases like MySQL)
	if b.limit > 0 {
		query.WriteString(fmt.Sprintf(" limit %d", b.limit))
	}

	return Query{
		SQL:    query.String(),
		Params: params,
	}
}

func (b *QueryBuilder) buildWhereClause(paramCount int) (string, []interface{}, int) {
	var query strings.Builder
	var params []interface{}

	query.WriteString(" where ")
	for i, where := range b.whereClauses {
		if i > 0 {
			query.WriteString(" " + where.JoinType + " ")
		}
		paramCount++
		query.WriteString(where.Column)
		query.WriteString(" " + where.Operator + " " + b.getPlaceholder(paramCount))
		params = append(params, where.Value)
	}

	return query.String(), params, paramCount
}

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

type Query struct {
	SQL    string
	Params []interface{}
}

type QueryBuilder struct {
	table        string
	columns      []string
	whereClauses []*WhereClause
	order        string
	limit        int
	offset       int
	paramStyle   ParameterStyle
}

type WhereClause struct {
	Column   string
	Operator string
	Value    interface{}
	JoinType string // AND/OR
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		columns:    []string{"*"},
		paramStyle: DollarNumber, // Default to DollarNumber
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

func (b *QueryBuilder) Select(columns ...string) *QueryBuilder {
	if len(columns) > 0 {
		b.columns = columns
	}
	return b
}

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

func (b *QueryBuilder) OrderBy(order string) *QueryBuilder {
	b.order = order
	return b
}

func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
	b.limit = limit
	return b
}

func (b *QueryBuilder) Offset(offset int) *QueryBuilder {
	b.offset = offset
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
	var query strings.Builder
	var params []interface{}
	paramCount := 0

	// Build SELECT clause
	query.WriteString("select ")
	query.WriteString(strings.Join(b.columns, ", "))

	// Build FROM clause
	query.WriteString(" from ")
	query.WriteString(b.table)

	// Build WHERE clause
	if len(b.whereClauses) > 0 {
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


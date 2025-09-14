package main

import (
	"fmt"
	"github.com/scape-labs/query"
)

func main() {
	// Test 1: Simple query without joins
	fmt.Println("Test 1: Simple query without joins")
	qb1 := query.NewQueryBuilder().
		Table("users").
		Select("id", "name", "email").
		Where("age", ">", 18)
		
	query1 := qb1.Build()
	fmt.Printf("SQL: %s\n", query1.Sql())
	fmt.Printf("Params: %v\n\n", query1.Params)

	// Test 2: Query with JOIN
	fmt.Println("Test 2: Query with JOIN")
	qb2 := query.NewQueryBuilder().
		Table("users").
		Select("users.id", "users.name", "accounts.name as account_name").
		LeftJoin("accounts", "accounts.id = users.account_id").
		Where("users.active", "=", true)
		
	query2 := qb2.Build()
	fmt.Printf("SQL: %s\n", query2.Sql())
	fmt.Printf("Params: %v\n\n", query2.Params)

	// Test 3: Query with multiple JOINs and aliases
	fmt.Println("Test 3: Query with multiple JOINs and aliases")
	qb3 := query.NewQueryBuilder().
		Table("posts").
		As("p").
		Select("p.title", "u.name as author", "c.name as category").
		LeftJoinAs("users", "u", "u.id = p.user_id").
		InnerJoinAs("categories", "c", "c.id = p.category_id").
		Where("p.published", "=", true).
		OrderBy("p.created_at").
		Limit(10)
		
	query3 := qb3.Build()
	fmt.Printf("SQL: %s\n", query3.Sql())
	fmt.Printf("Params: %v\n\n", query3.Params)

	fmt.Println("All tests passed!")
}
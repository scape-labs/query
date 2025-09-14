package query

import (
	"strings"
	"testing"
)

// Basic Query Tests

func TestBasicSelectQuery(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name", "email")

	query := qb.Build()
	expectedSQL := "select id, name, email from users"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 0 {
		t.Errorf("Expected no params, got: %v", query.Params)
	}
}

func TestSelectQueryWithWhereClause(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name", "email").
		Where("age", ">", 18).
		Where("active", "=", true)

	query := qb.Build()
	expectedSQL := "select id, name, email from users where age > $1 and active = $2"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 2 || query.Params[0] != 18 || query.Params[1] != true {
		t.Errorf("Expected params: [18, true], got: %v", query.Params)
	}
}

func TestSelectQueryWithOrWhereClause(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name", "email").
		Where("age", ">", 18).
		OrWhere("admin", "=", true)

	query := qb.Build()
	expectedSQL := "select id, name, email from users where age > $1 or admin = $2"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 2 || query.Params[0] != 18 || query.Params[1] != true {
		t.Errorf("Expected params: [18, true], got: %v", query.Params)
	}
}

func TestSelectQueryWithOrderByLimitOffset(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name", "email").
		OrderBy("name").
		Limit(10).
		Offset(20)

	query := qb.Build()
	expectedSQL := "select id, name, email from users order by name limit 10 offset 20"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 0 {
		t.Errorf("Expected no params, got: %v", query.Params)
	}
}

func TestSelectQueryWithAllClauses(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name", "email").
		Where("age", ">", 18).
		Where("active", "=", true).
		OrderBy("name").
		Limit(10).
		Offset(20)

	query := qb.Build()
	expectedSQL := "select id, name, email from users where age > $1 and active = $2 order by name limit 10 offset 20"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 2 || query.Params[0] != 18 || query.Params[1] != true {
		t.Errorf("Expected params: [18, true], got: %v", query.Params)
	}
}

func TestSelectQueryWithWildcard(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users")

	query := qb.Build()
	expectedSQL := "select * from users"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 0 {
		t.Errorf("Expected no params, got: %v", query.Params)
	}
}

// INSERT Query Tests

func TestBasicInsertQuery(t *testing.T) {
	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}

	qb := NewQueryBuilder().
		Table("users").
		Insert(data)

	query := qb.Build()
	
	// Check that SQL contains the expected components (order may vary due to map iteration)
	if !strings.Contains(query.SQL, "insert into users") {
		t.Errorf("Expected SQL to contain 'insert into users', got: %s", query.SQL)
	}
	
	if !strings.Contains(query.SQL, "(name, email, age)") && !strings.Contains(query.SQL, "(age, email, name)") &&
		!strings.Contains(query.SQL, "(email, name, age)") && !strings.Contains(query.SQL, "(name, age, email)") &&
		!strings.Contains(query.SQL, "(email, age, name)") && !strings.Contains(query.SQL, "(age, name, email)") {
		t.Errorf("Expected SQL to contain column list with name, email, age, got: %s", query.SQL)
	}
	
	if !strings.Contains(query.SQL, "values ($1, $2, $3)") {
		t.Errorf("Expected SQL to contain 'values ($1, $2, $3)', got: %s", query.SQL)
	}

	if len(query.Params) != 3 {
		t.Errorf("Expected 3 params, got: %v", query.Params)
	}

	// Check that params match expected values (order may vary due to map iteration)
	params := make(map[interface{}]bool)
	for _, param := range query.Params {
		params[param] = true
	}

	if !params[30] || !params["john@example.com"] || !params["John Doe"] {
		t.Errorf("Expected params to contain 30, 'john@example.com', and 'John Doe', got: %v", query.Params)
	}
}

func TestInsertQueryWithSpecificColumns(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		InsertColumns("name", "email").
		Values("John Doe", "john@example.com")

	query := qb.Build()
	expectedSQL := "insert into users (name, email) values ($1, $2)"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 2 || query.Params[0] != "John Doe" || query.Params[1] != "john@example.com" {
		t.Errorf("Expected params: ['John Doe', 'john@example.com'], got: %v", query.Params)
	}
}

// UPDATE Query Tests

func TestBasicUpdateQuery(t *testing.T) {
	data := map[string]interface{}{
		"name":  "Jane Doe",
		"email": "jane@example.com",
	}

	qb := NewQueryBuilder().
		Table("users").
		Update(data)

	query := qb.Build()
	
	// Check that SQL contains the expected components (order may vary due to map iteration)
	if !strings.Contains(query.SQL, "update users set") {
		t.Errorf("Expected SQL to contain 'update users set', got: %s", query.SQL)
	}
	
	if !strings.Contains(query.SQL, "email = $1") && !strings.Contains(query.SQL, "name = $1") {
		t.Errorf("Expected SQL to contain column assignments, got: %s", query.SQL)
	}

	if len(query.Params) != 2 {
		t.Errorf("Expected 2 params, got: %v", query.Params)
	}

	// Check that params match expected values (order may vary due to map iteration)
	params := make(map[interface{}]bool)
	for _, param := range query.Params {
		params[param] = true
	}

	if !params["jane@example.com"] || !params["Jane Doe"] {
		t.Errorf("Expected params to contain 'jane@example.com' and 'Jane Doe', got: %v", query.Params)
	}
}

func TestUpdateQueryWithWhereClause(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Update(map[string]interface{}{
			"name":  "Jane Doe",
			"email": "jane@example.com",
		}).
		Where("id", "=", 1)

	query := qb.Build()
	
	// Check that SQL contains the expected components (order may vary due to map iteration)
	if !strings.Contains(query.SQL, "update users set") {
		t.Errorf("Expected SQL to contain 'update users set', got: %s", query.SQL)
	}
	
	if !strings.Contains(query.SQL, "email = $1") && !strings.Contains(query.SQL, "name = $1") {
		t.Errorf("Expected SQL to contain column assignments, got: %s", query.SQL)
	}
	
	if !strings.Contains(query.SQL, "where id = $3") && !strings.Contains(query.SQL, "where id = $2") {
		t.Errorf("Expected SQL to contain where clause, got: %s", query.SQL)
	}

	if len(query.Params) != 3 {
		t.Errorf("Expected 3 params, got: %v", query.Params)
	}

	// Check that params match expected values (order may vary due to map iteration)
	params := make(map[interface{}]bool)
	for _, param := range query.Params {
		params[param] = true
	}

	if !params["jane@example.com"] || !params["Jane Doe"] || !params[1] {
		t.Errorf("Expected params to contain 'jane@example.com', 'Jane Doe', and 1, got: %v", query.Params)
	}
}

func TestUpdateQueryWithSetMethod(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Set("name", "Jane Doe").
		Set("email", "jane@example.com").
		Where("id", "=", 1)

	query := qb.Build()
	expectedSQL := "update users set name = $1, email = $2 where id = $3"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 3 || query.Params[0] != "Jane Doe" || query.Params[1] != "jane@example.com" || query.Params[2] != 1 {
		t.Errorf("Expected params: ['Jane Doe', 'jane@example.com', 1], got: %v", query.Params)
	}
}

// DELETE Query Tests

func TestBasicDeleteQuery(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Delete()

	query := qb.Build()
	expectedSQL := "delete from users"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 0 {
		t.Errorf("Expected no params, got: %v", query.Params)
	}
}

func TestDeleteQueryWithWhereClause(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Delete().
		Where("id", "=", 1)

	query := qb.Build()
	expectedSQL := "delete from users where id = $1"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 1 || query.Params[0] != 1 {
		t.Errorf("Expected params: [1], got: %v", query.Params)
	}
}

func TestDeleteQueryWithMultipleWhereClauses(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Delete().
		Where("id", "=", 1).
		Where("active", "=", false)

	query := qb.Build()
	expectedSQL := "delete from users where id = $1 and active = $2"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 2 || query.Params[0] != 1 || query.Params[1] != false {
		t.Errorf("Expected params: [1, false], got: %v", query.Params)
	}
}

// Parameter Style Tests

func TestQuestionMarkParameterStyle(t *testing.T) {
	qb := NewQueryBuilder().
		ParameterPlaceholder(QuestionMark).
		Table("users").
		Select("id", "name").
		Where("age", ">", 18)

	query := qb.Build()
	expectedSQL := "select id, name from users where age > ?"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 1 || query.Params[0] != 18 {
		t.Errorf("Expected params: [18], got: %v", query.Params)
	}
}

func TestDollarNumberParameterStyle(t *testing.T) {
	qb := NewQueryBuilder().
		ParameterPlaceholder(DollarNumber).
		Table("users").
		Select("id", "name").
		Where("age", ">", 18)

	query := qb.Build()
	expectedSQL := "select id, name from users where age > $1"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 1 || query.Params[0] != 18 {
		t.Errorf("Expected params: [18], got: %v", query.Params)
	}
}

// JOIN Operation Tests

func TestJoinOperations(t *testing.T) {
	// Test basic JOIN
	qb := NewQueryBuilder().
		Table("users").
		Select("users.id", "users.name", "accounts.name as account_name").
		Join("accounts", "accounts.id = users.account_id").
		Where("users.active", "=", true)

	query := qb.Build()
	expectedSQL := "select users.id, users.name, accounts.name as account_name from users JOIN accounts on accounts.id = users.account_id where users.active = $1"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 1 || query.Params[0] != true {
		t.Errorf("Expected params: [true], got: %v", query.Params)
	}
}

func TestLeftJoinOperation(t *testing.T) {
	qb := NewQueryBuilder().
		Table("posts").
		As("p").
		Select("p.title", "u.name as author").
		LeftJoinAs("users", "u", "u.id = p.user_id").
		Where("p.published", "=", true)

	query := qb.Build()
	expectedSQL := "select p.title, u.name as author from posts as p LEFT JOIN users as u on u.id = p.user_id where p.published = $1"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 1 || query.Params[0] != true {
		t.Errorf("Expected params: [true], got: %v", query.Params)
	}
}

func TestMultipleJoins(t *testing.T) {
	qb := NewQueryBuilder().
		Table("orders").
		Select("orders.id", "customers.name", "products.name as product_name").
		LeftJoin("customers", "customers.id = orders.customer_id").
		InnerJoin("order_items", "order_items.order_id = orders.id").
		LeftJoin("products", "products.id = order_items.product_id").
		Where("orders.status", "=", "completed")

	query := qb.Build()
	expectedSQL := "select orders.id, customers.name, products.name as product_name from orders LEFT JOIN customers on customers.id = orders.customer_id INNER JOIN order_items on order_items.order_id = orders.id LEFT JOIN products on products.id = order_items.product_id where orders.status = $1"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 1 || query.Params[0] != "completed" {
		t.Errorf("Expected params: [completed], got: %v", query.Params)
	}
}

func TestFullJoinOperations(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("users.name", "accounts.name as account_name").
		FullJoin("accounts", "accounts.id = users.account_id")

	query := qb.Build()
	expectedSQL := "select users.name, accounts.name as account_name from users FULL JOIN accounts on accounts.id = users.account_id"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 0 {
		t.Errorf("Expected no params, got: %v", query.Params)
	}
}

func TestRightJoinOperations(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("users.name", "accounts.name as account_name").
		RightJoin("accounts", "accounts.id = users.account_id")

	query := qb.Build()
	expectedSQL := "select users.name, accounts.name as account_name from users RIGHT JOIN accounts on accounts.id = users.account_id"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 0 {
		t.Errorf("Expected no params, got: %v", query.Params)
	}
}

func TestInnerJoinOperations(t *testing.T) {
	qb := NewQueryBuilder().
		Table("users").
		Select("users.name", "accounts.name as account_name").
		InnerJoin("accounts", "accounts.id = users.account_id")

	query := qb.Build()
	expectedSQL := "select users.name, accounts.name as account_name from users INNER JOIN accounts on accounts.id = users.account_id"
	if query.SQL != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, query.SQL)
	}

	if len(query.Params) != 0 {
		t.Errorf("Expected no params, got: %v", query.Params)
	}
}

// Benchmark Tests

func BenchmarkSelectQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qb := NewQueryBuilder().
			Table("users").
			Select("id", "name", "email").
			Where("age", ">", 18).
			Where("active", "=", true).
			OrderBy("name").
			Limit(10)

		_ = qb.Build()
	}
}

func BenchmarkInsertQuery(b *testing.B) {
	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}

	for i := 0; i < b.N; i++ {
		qb := NewQueryBuilder().
			Table("users").
			Insert(data)

		_ = qb.Build()
	}
}

func BenchmarkUpdateQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qb := NewQueryBuilder().
			Table("users").
			Update(map[string]interface{}{
				"name":  "Jane Doe",
				"email": "jane@example.com",
			}).
			Where("id", "=", 1)

		_ = qb.Build()
	}
}

func BenchmarkDeleteQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qb := NewQueryBuilder().
			Table("users").
			Delete().
			Where("id", "=", 1)

		_ = qb.Build()
	}
}

func BenchmarkJoinQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qb := NewQueryBuilder().
			Table("users").
			Select("users.id", "users.name", "accounts.name as account_name").
			LeftJoin("accounts", "accounts.id = users.account_id").
			Where("users.active", "=", true)

		_ = qb.Build()
	}
}
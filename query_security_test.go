package query

import (
	"strings"
	"testing"
)

// SQL Injection Tests
func TestSQLInjection_TableName(t *testing.T) {
	// This test demonstrates the vulnerability - table name is not properly escaped
	qb := NewQueryBuilder().
		Table("users; DROP TABLE accounts; --").
		Select("id", "name")

	query := qb.Build()
	
	// This should NOT contain the malicious SQL
	if strings.Contains(query.SQL, "DROP TABLE accounts") {
		t.Errorf("SQL injection vulnerability detected in table name: %s", query.SQL)
	}
	
	// The table name should be properly escaped or quoted
	// This is currently failing, showing the vulnerability exists
}

func TestSQLInjection_ColumnName(t *testing.T) {
	// This test demonstrates the vulnerability - column name is not properly escaped
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name, email; DROP TABLE accounts; --")

	query := qb.Build()
	
	// This should NOT contain the malicious SQL
	if strings.Contains(query.SQL, "DROP TABLE accounts") {
		t.Errorf("SQL injection vulnerability detected in column name: %s", query.SQL)
	}
}

func TestSQLInjection_WhereColumn(t *testing.T) {
	// This test demonstrates the vulnerability - WHERE column is not properly escaped
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name").
		Where("id; DROP TABLE accounts; --", "=", 1)

	query := qb.Build()
	
	// This should NOT contain the malicious SQL
	if strings.Contains(query.SQL, "DROP TABLE accounts") {
		t.Errorf("SQL injection vulnerability detected in WHERE column: %s", query.SQL)
	}
}

func TestSQLInjection_OrderBy(t *testing.T) {
	// This test demonstrates the vulnerability - ORDER BY is not properly escaped
	qb := NewQueryBuilder().
		Table("users").
		Select("id", "name").
		OrderBy("name; DROP TABLE accounts; --")

	query := qb.Build()
	
	// This should NOT contain the malicious SQL
	if strings.Contains(query.SQL, "DROP TABLE accounts") {
		t.Errorf("SQL injection vulnerability detected in ORDER BY: %s", query.SQL)
	}
}

func TestSQLInjection_JoinTable(t *testing.T) {
	// This test demonstrates the vulnerability - JOIN table is not properly escaped
	qb := NewQueryBuilder().
		Table("users").
		Select("users.id", "users.name", "accounts.name as account_name").
		LeftJoin("accounts; DROP TABLE logs; --", "accounts.id = users.account_id")

	query := qb.Build()
	
	// This should NOT contain the malicious SQL
	if strings.Contains(query.SQL, "DROP TABLE logs") {
		t.Errorf("SQL injection vulnerability detected in JOIN table: %s", query.SQL)
	}
}

func TestSQLInjection_JoinCondition(t *testing.T) {
	// This test demonstrates the vulnerability - JOIN condition is not properly escaped
	qb := NewQueryBuilder().
		Table("users").
		Select("users.id", "users.name", "accounts.name as account_name").
		LeftJoin("accounts", "accounts.id = users.account_id; DROP TABLE logs; --")

	query := qb.Build()
	
	// This should NOT contain the malicious SQL
	if strings.Contains(query.SQL, "DROP TABLE logs") {
		t.Errorf("SQL injection vulnerability detected in JOIN condition: %s", query.SQL)
	}
}
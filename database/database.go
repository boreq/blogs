// Package database provides database access.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/boreq/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type DatabaseType int

const (
	PostgreSQL DatabaseType = iota
	SQLite3
)

// DB becomes initialized after calling Init.
var DB *sqlx.DB

// ErrNoRows lets the user access sql.ErrNoRows without importing database/sql.
var ErrNoRows = sql.ErrNoRows

var createTableQueries = []string{
	createUserSQL,
	createUserSessionSQL,
	createBlogSQL,
	createCategorySQL,
	createPostSQL,
	createTagSQL,
	createPostToTagSQL,
	createUpdateSQL,
	createSubscriptionSQL,
	createStarSQL,
	createInsertStarTriggerSQL,
	createDeleteStarTriggerSQL,
	createInsertSubscriptionTriggerSQL,
	createDeleteSubscriptionTriggerSQL,
}

var tableNames = []string{
	"user_session",
	"user",
	"post",
	"post_to_tag",
	"tag",
	"update",
	"category",
	"blog",
	"subscription",
	"star",
}

// Init connects to the specified database.
func Init(databaseType DatabaseType, params string) (err error) {
	switch databaseType {
	case SQLite3:
		DB, err = sqlx.Connect("sqlite3", params)
		if err != nil {
			return err
		}
		break
	case PostgreSQL:
		DB, err = sqlx.Connect("postgres", params)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("Reached the default switch case in database.Init")
	}
	DB.MapperFunc(mapperFunc)
	return nil
}

// CreateTables creates database tables.
func CreateTables() error {
	for _, query := range createTableQueries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

// DropTables drops all database tables used by this program.
func DropTables() error {
	for _, tableName := range tableNames {
		query := fmt.Sprintf("DROP TABLE IF EXISTS \"%s\"", tableName)
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func mapperFunc(fieldName string) string {
	var result string
	for i, ch := range fieldName {
		if i > 0 && i < len(fieldName)-1 && ch > 'A' && ch < 'Z' {
			result += "_"
		}
		result += strings.ToLower(string(ch))
	}
	return result
}

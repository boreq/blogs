package database

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseType int

const (
	PostgreSQL DatabaseType = iota
	SQLite3
)

// DB becomes initialized after calling Init.
var DB *gorm.DB

// Init connects to the specified database.
// http://jinzhu.me/gorm/database.html#connecting-to-a-database
func Init(databaseType DatabaseType, params string) (err error) {
	switch databaseType {
	case SQLite3:
		DB, err = gorm.Open("sqlite3", params)
		break
	case PostgreSQL:
		DB, err = gorm.Open("postgres", params)
		break
	default:
		return errors.New("Reached the default switch case in database.Init")
	}
	return err
}

var tables = []interface{}{&Blog{}, &Category{}, &Post{}, &Tag{}}

// MigrateTables creates missing tables and columns.
func MigrateTables() error {
	DB.AutoMigrate(tables...)
	return nil
}

// DropTables drops all tables used by this program.
func DropTables() error {
	DB.DropTable(tables...)
	return nil
}

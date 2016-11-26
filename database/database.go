package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseType int

const (
	PostgreSQL DatabaseType = iota
	SQLite
)

var DB *gorm.DB

func Init(databaseType DatabaseType, params string) error {
	var err error
	DB, err = gorm.Open("sqlite3", params)
	if err != nil {
		return err
	}
	return nil
}

var tables = []interface{}{&Blog{}, &Category{}, &Post{}, &Tag{}}

func MigrateTables() error {
	DB.AutoMigrate(tables...)
	return nil
}

func DropTables() error {
	DB.DropTable(tables...)
	return nil
}

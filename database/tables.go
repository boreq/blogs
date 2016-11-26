package database

import (
	"time"
)

type Blog struct {
	ID uint `gorm:"primary_key"`

	InternalID uint   `gorm:"not null"`
	Title      string `gorm:"not null"`
	Categories []Category
}

type Category struct {
	ID     uint `gorm:"primary_key"`
	BlogID uint

	Name  string `gorm:"not null"`
	Posts []Post
}

type Post struct {
	ID         uint `gorm:"primary_key"`
	CategoryID uint

	InternalID string `gorm:"not null"`
	Title      string `gorm:"not null"`
	Date       time.Time
	Tags       []Tag `gorm:"many2many:post_tags;"`
}

type Tag struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null"`
}

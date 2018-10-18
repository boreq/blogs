package dto

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/utils"
	"github.com/pkg/errors"
	"time"
)

type ScannableTime struct {
	time.Time
}

func (t *ScannableTime) Scan(src interface{}) error {
	switch src := src.(type) {
	case time.Time:
		t.Time = src
	case []uint8:
		tmp, err := time.Parse("2006-01-02 15:04:05-07:00", string(src))
		if err != nil {
			return err
		}
		t.Time = tmp
	default:
		return errors.New("Invalid type in Scan")
	}
	return nil
}

func (t ScannableTime) String() string {
	return utils.ISO8601(t.Time)
}

type Page struct {
	// Page is the page number.
	Page int `json:"page"`

	// PerPage specifies the number of items on a single page.
	PerPage int `json:"perPage"`
}

type PageOut struct {
	Page

	// AllItems is the number of all available items.
	AllItems int `json:"allItems"`
}

type BlogOut struct {
	database.Blog
	LastPost   *ScannableTime `json:"lastPost,omitempty"`
	Url        string         `json:"url"`
	CleanUrl   string         `json:"cleanUrl"`
	Subscribed bool           `json:"subscribed"`
}

type PostOut struct {
	database.Post
	Category database.Category `json:"category"`
	Blog     database.Blog     `json:"blog"`
	Url      string            `json:"url"`
	Starred  bool              `json:"starred"`
}

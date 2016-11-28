package updater

import (
	"encoding/json"
	"github.com/boreq/blogs/database"
)

type stats struct {
	TitleDownloaded bool
	TitleUpdated    bool
	TitleCorrect    bool

	LoaderErrors int

	PostsReceived  int
	PostsUpdated   int
	PostsUnaltered int

	TagAddErrors    int
	TagRemoveErrors int

	PostRemovalsStarted   bool
	PostRemovalsAttempted int
	PostRemovalsSucceeded int
}

type update struct {
	database.Update
	stats
}

func (u update) GetSucceeded() bool {
	if !(u.TitleDownloaded && (u.TitleUpdated || u.TitleCorrect)) {
		return false
	}
	if u.LoaderErrors > 0 {
		return false
	}
	if u.PostsReceived != u.PostsUpdated+u.PostsUnaltered {
		return false
	}
	if u.TagAddErrors > 0 || u.TagRemoveErrors > 0 {
		return false
	}
	if u.LoaderErrors == 0 && !u.PostRemovalsStarted {
		return false
	}
	if u.PostRemovalsAttempted != u.PostRemovalsSucceeded {
		return false
	}
	return true
}

func (u update) GetData() (string, error) {
	b, err := json.Marshal(u.stats)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func saveStatistics(u update) error {
	json, err := u.GetData()
	if err != nil {
		return err
	}
	_, err = database.DB.Exec("INSERT INTO \"update\" (blog_id, started, ended, succeeded, data) VALUES ($1, $2, $3, $4, $5)",
		u.Update.BlogID, u.Started, u.Ended, u.GetSucceeded(), json)
	return err
}

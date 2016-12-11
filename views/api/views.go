package api

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/http/api"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sort"
)

type updatesChartResponse struct {
	Labels  []string `json:"labels"`
	Success []int    `json:"success"`
	Failure []int    `json:"failure"`
}

type entry struct {
	Label   string
	Success int
	Failure int
}

type byLabel []entry

func (d byLabel) Len() int           { return len(d) }
func (d byLabel) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d byLabel) Less(i, j int) bool { return d[i].Label < d[j].Label }

func updatesChart(r *http.Request, p httprouter.Params) (interface{}, api.Error) {
	var updates []database.Update
	if err := database.DB.Select(&updates,
		`SELECT "update".*
		FROM "update"
		WHERE started > (SELECT DATETIME('now', '-30 day'))`); err != nil {
		return nil, api.InternalServerError
	}

	// This is a nightmare to do if you need to support multiple databases
	tmp := make(map[string]*entry)
	for _, update := range updates {
		label := update.Started.Format("2006-01-02")
		if _, ok := tmp[label]; !ok {
			tmp[label] = &entry{Label: label}
		}
		if update.Succeeded {
			tmp[label].Success++
		} else {
			tmp[label].Failure++
		}
	}

	var sorted []entry
	for _, v := range tmp {
		sorted = append(sorted, *v)
	}
	sort.Sort(byLabel(sorted))

	response := updatesChartResponse{}
	for _, e := range sorted {
		response.Labels = append(response.Labels, e.Label)
		response.Success = append(response.Success, e.Success)
		response.Failure = append(response.Failure, e.Failure)
	}
	return response, nil
}

package utils

import (
	"net/http"
	"strings"
)

type SortParam struct {
	Key      string
	Label    string
	Query    string
	Reversed bool
}

type SortCriteria struct {
	Label        string
	Selected     bool
	key          string
	selectedDesc bool
	reversed     bool
}

func (c SortCriteria) GetKey() string {
	if c.shouldSortDesc() {
		return c.key + "_desc"
	} else {
		return c.key
	}
}

func (c SortCriteria) shouldSortDesc() bool {
	if c.Selected {
		if c.selectedDesc {
			return false
		} else {
			return true
		}
	} else {
		if c.reversed {
			return true
		} else {
			return false
		}
	}
}

func (c SortCriteria) GetCurrentKey() string {
	if c.shouldSortDesc() {
		return c.key
	} else {
		return c.key + "_desc"
	}
}

func (c SortCriteria) IsAsc() bool {
	return c.shouldSortDesc()
}

type Sort struct {
	// Query can be directly used in the ORDER BY clause of an SQL query.
	Query string

	// Criteria can be used to generate a list in the templates.
	Criteria []SortCriteria

	// Current key is the value of the "sort" query parameter or the
	// default value if the parameter is invalid. A parameter is invalid
	// if it doesn't exist in the parameter list passed to the NewSort.
	// The default value is the key of the first element of the parameter
	// list passed to NewSort.
	CurrentKey string
}

// NewSort uses the "sort" query parameter to get the sort key and initialize
// the struct.
func NewSort(r *http.Request, params []SortParam) Sort {
	sortParams, ok := r.URL.Query()["sort"]
	var sortParam *SortParam
	var desc bool
	if ok {
		sortParam = getSortParam(sortParams[0], params)
		desc = strings.HasSuffix(sortParams[0], "_desc")
	}
	if sortParam == nil {
		sortParam = &params[0]
		desc = sortParam.Reversed
	}

	query := sortParam.Query
	if desc {
		query += " DESC"
	}

	currentKey := sortParam.Key
	if desc {
		currentKey += "_desc"
	}

	rv := Sort{
		Query:      query,
		CurrentKey: currentKey,
	}
	for _, param := range params {
		rv.Criteria = append(rv.Criteria, SortCriteria{
			key:          param.Key,
			Label:        param.Label,
			Selected:     param == *sortParam,
			selectedDesc: param == *sortParam && desc,
			reversed:     param.Reversed,
		})
	}
	return rv
}

func getSortParam(key string, sortParams []SortParam) *SortParam {
	key = strings.TrimSuffix(key, "_desc")
	for _, sortParam := range sortParams {
		if key == sortParam.Key {
			return &sortParam
		}
	}
	return nil
}

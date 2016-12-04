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
	if c.Selected && !c.selectedDesc {
		return c.key + "_desc"
	}
	return c.key
}

func (c SortCriteria) IsAsc() bool {
	if c.Selected {
		if c.reversed {
			return c.selectedDesc
		} else {
			return !c.selectedDesc
		}
	} else {
		if c.reversed {
			return true
		} else {
			return false
		}
	}
}

type Sort struct {
	Query    string
	Criteria []SortCriteria
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
	}

	query := sortParam.Query
	if (desc && !sortParam.Reversed) || (!desc && sortParam.Reversed) {
		query += " DESC"
	}

	rv := Sort{Query: query}
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

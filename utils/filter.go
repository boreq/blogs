package utils

import (
	"net/http"
)

type FilterParam struct {
	Key   string
	Label string
	Query string
}

type FilterCriteria struct {
	Key      string
	Label    string
	Selected bool
}

type Filter struct {
	// Query can be used in the WHERE clause of an SQL query.
	Query string

	// Criteria can be used to generate a list in the templates.
	Criteria []FilterCriteria

	// Current key is the value of the "filter" query parameter or the
	// default value if the parameter is invalid. A parameter is invalid
	// if it doesn't exist in the parameter list passed to the NewFilter.
	// The default value is the key of the first element of the parameter
	// list passed to NewFilter.
	CurrentKey string
}

// NewFilter uses the "filter" query parameter to get the filter key and
// initialize the struct.
func NewFilter(r *http.Request, params []FilterParam) Filter {
	var param *FilterParam
	if queryParams, ok := r.URL.Query()["filter"]; ok {
		param = getFilterParam(queryParams[0], params)
	}
	if param == nil {
		param = &params[0]
	}

	rv := Filter{
		Query:      param.Query,
		CurrentKey: param.Key,
	}
	for _, p := range params {
		rv.Criteria = append(rv.Criteria, FilterCriteria{
			Key:      p.Key,
			Label:    p.Label,
			Selected: p == *param,
		})
	}
	return rv
}

func getFilterParam(key string, filterParams []FilterParam) *FilterParam {
	for _, filterParam := range filterParams {
		if key == filterParam.Key {
			return &filterParam
		}
	}
	return nil
}

package templates

import (
	"fmt"
	"html/template"
)

func funcDict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict accepts an even number of parameters, %d given", len(values))
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func funcMinus(a, b int) int {
	return a - b
}

func funcPlus(a, b int) int {
	return a + b
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict":  funcDict,
		"minus": funcMinus,
		"plus":  funcPlus,
	}
}

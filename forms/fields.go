package forms

import (
	"fmt"
	"html/template"
)

type TextField struct {
	Name       string
	Label      string
	HelpText   string
	Cleaners   []Cleaner
	Validators []Validator

	attributes map[string]string
	errors     []string
	value      string
}

func (f *TextField) SetAttribute(key, value string) *TextField {
	if f.attributes == nil {
		f.attributes = make(map[string]string)
	}
	f.attributes[key] = value
	return f
}

func (f TextField) Render() template.HTML {
	if f.attributes == nil {
		f.attributes = make(map[string]string)
	}
	if _, ok := f.attributes["type"]; !ok {
		f.attributes["type"] = "text"
	}
	if _, ok := f.attributes["id"]; !ok {
		f.attributes["id"] = "id_" + f.Name
	}
	f.attributes["name"] = f.GetName()
	f.attributes["value"] = f.GetValue()
	attributes := ""
	for key, value := range f.attributes {
		attributes += fmt.Sprintf(" %s=\"%s\"", key, value)
	}
	str := "<input" + attributes + ">"
	return template.HTML(str)
}

func (f TextField) GetName() string {
	return f.Name
}

func (f TextField) GetLabel() string {
	return f.Label
}

func (f TextField) GetHelpText() string {
	return f.HelpText
}

func (f *TextField) Validate(value string) bool {
	for _, cleaner := range f.Cleaners {
		value = cleaner(value)
	}
	f.value = value
	for _, validator := range f.Validators {
		f.errors = append(f.errors, validator(value)...)
	}
	return f.IsValid()
}

func (f *TextField) AddError(err string) {
	f.errors = append(f.errors, err)
}

func (f TextField) Errors() ValidationErrors {
	return f.errors
}

func (f TextField) IsValid() bool {
	return f.errors == nil || len(f.errors) == 0
}

func (f TextField) GetValue() string {
	return f.value
}

func ToPasswordField(f *TextField) *TextField {
	f.SetAttribute("type", "password")
	return f
}

package main

import (
	"bytes"
	"html/template"
)

var tmpl = template.New("TemplateContent")

func ParseContent(content string, param interface{}) (string, error) {
	if len(content) == 0 {
		return "", nil
	}
	newTmpl, err := tmpl.Clone()
	if err != nil {
		return "", err
	}
	newTmpl, err = newTmpl.Parse(content)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBuffer(nil)
	err = newTmpl.Execute(buffer, param)
	if err != nil {
		return "", err
	}
	return buffer.String(), err
}

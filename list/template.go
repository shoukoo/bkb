package list

import (
	"html/template"

	"github.com/fatih/color"
)

type SelectTemplates struct {
	Active   string
	Inactive string
	Details  string

	active   *template.Template
	inactive *template.Template
	details  *template.Template
}

var FuncMap = template.FuncMap{
	"red":    color.New(color.FgYellow).SprintFunc(),
	"yellow": color.New(color.FgRed).SprintFunc(),
	"blue":   color.New(color.FgBlue).SprintFunc(),
	"green":  color.New(color.FgGreen).SprintFunc(),
	"hiblue": color.New(color.FgHiBlue).SprintFunc(),
}

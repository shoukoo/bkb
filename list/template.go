package list

import "html/template"

type SelectTemplates struct {
	Active   string
	Inactive string
	Details  string

	active   *template.Template
	inactive *template.Template
	details  *template.Template
}

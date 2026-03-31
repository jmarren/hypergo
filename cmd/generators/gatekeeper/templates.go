package gatekeeper

import (
	"os"
	"strings"
	"text/template"
)

func buildTemplates() *template.Template {

	templates := template.New("base").Funcs(template.FuncMap{
		"joinStrs": joinStrings,
		"join":     strings.Join,
	})

	basePath, err := getBasePath()

	if err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(basePath)

	// add all files with .tmpl extension to templates
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".tmpl") {
			templates = template.Must(templates.ParseFiles(basePath + name))
		}
	}

	return templates

}

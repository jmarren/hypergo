package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of gatekeeper:\n")
	fmt.Fprintf(os.Stderr, "\tgatekeeper [path to yaml file]\n")
	flag.PrintDefaults()
}

type Field struct {
	Name string
	Kind string
	// Required bool
}

type IntField struct {
	Field
	Min int
	Max int
}

type StringField struct {
	Field
	MaxLen int
	MinLen int
}

type Object struct {
	Name   string
	Fields []Field
}

type Config struct {
	Package string
	Objects []Object
}

func (f *Field) MakeValidator() {
	tmpl, err := template.New("validator").Parse(`
{{  if eq .Kind  "int" }}_, err := strconv.Atoi({{ .Name }}) { 
	if err != nil {
		errs = 	append(errs, "not of type int") 
	}
}{{ end }}
`)

	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, f)

	if err != nil {
		panic(err)
	}
}

func (f *Field) intValidation() {

	tmpl, err := template.New("intValidation").Parse(`
			x.{{ .Name }}, err = strconv.Atoi(r.FormValue("{{ .Name }}"))
        if err != nil {
		errs = append(errs, fmt.Errorf("{{ .Name }} must be a number"))
        }

	`)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, f)
	if err != nil {
		panic(err)
	}

}

func (f *Field) Validation() {

	switch f.Kind {
	case "int":
		f.intValidation()
	}
}

// TemplateRegistry is a custom HTML template renderer for Echo framework
type TemplateRegistry struct {
	templates *template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var TmplRegistry *TemplateRegistry

// Get current file's absolute path
func getBasePath() (string, error) {
	_, file, _, ok := runtime.Caller(1) // Use caller(1) to get the path of the function that calls this helper
	if !ok {
		return "", fmt.Errorf("could not get caller information")
	}

	basePath, _ := strings.CutSuffix(file, "main.go")
	// runtime.Caller returns the path as known at compile time, use filepath.Abs for absolute path if needed
	// The path returned by runtime.Caller is already an absolute path in most cases.
	return basePath, nil
}

func main() {

	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	templates := template.New("base").Funcs(template.FuncMap{})

	basePath, err := getBasePath()

	if err != nil {
		panic(err)
	}

	// Parse base layout
	templates = template.Must(templates.ParseFiles(basePath + "base.tmpl"))
	templates = template.Must(templates.ParseFiles(basePath + "typedef.tmpl"))
	templates = template.Must(templates.ParseFiles(basePath + "constructor.tmpl"))
	templates = template.Must(templates.ParseFiles(basePath + "kind.tmpl"))
	templates = template.Must(templates.ParseFiles(basePath + "int.tmpl"))
	templates = template.Must(templates.ParseFiles(basePath + "string.tmpl"))

	// // Parse all partial templates (blocks)
	// for _, partial := range partials {
	// 	fmt.Println(partial)
	// 	templates = template.Must(templates.ParseFiles(dir + "/" + partial))
	// }

	// Set up the global template registry
	TmplRegistry = &TemplateRegistry{templates: templates}

	path := os.Args[1]

	fmt.Printf("path = %s\n", path)

	yamlBytes, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	t := Config{}

	yaml.Unmarshal(yamlBytes, &t)

	// 	tmpl, err := template.New("structBuilder").Parse(
	// 		`{{ define "header" }}
	// package {{ .Package }}
	//
	// import (
	// 	"net/http"
	// 	"strconv"
	// 	"fmt"
	// )
	// {{ end }}
	// {{ template "header" . }}
	//
	// {{ range .Objects }}
	// type {{ .Name }} struct {
	// {{ range .Fields }}	{{ .Name }} {{ .Kind }}
	// {{ end }}}
	//
	// func New{{ .Name }}(r *http.Request) (*{{ .Name }}, []error) {
	// 	errs := []error{}
	// 	var err error
	// 	x := new({{.Name }})
	// {{ range .Fields }}
	// {{ if eq .Kind "string" }}	x.{{ .Name }} = r.FormValue("{{ .Name }}"){{ end }}{{ if eq .Kind "int" }}	x.{{ .Name }}, err = strconv.Atoi(r.FormValue("{{ .Name }}"))
	//         if err != nil {
	// 		errs = append(errs, fmt.Errorf("{{ .Name }} must be a number"))
	//         }
	//       {{ end }}
	//       {{ end }}
	//
	//       return x, errs
	//
	// }
	//
	// {{ end }}
	//
	//
	//
	// `)

	// func (s *{{ .Name }}) FromRequest(r *http.Request) {
	//   s.{{ .Name }} = r.FormValue("{{ .Name }}")
	//
	// {{ end }}

	if err != nil {
		panic(err)
	}

	fmt.Printf("t.Objects[0] = %v\n", t.Objects[0].Fields[0])

	for _, field := range t.Objects[0].Fields {
		field.Validation()
	}

	file, err := os.OpenFile("gatekeeper.go", os.O_WRONLY|os.O_CREATE, 0777)

	err = templates.ExecuteTemplate(file, "base", t)
	if err != nil {
		panic(err)
	}
	//
	defer file.Close()

	// os.WriteFile()

	// fmt.Fprint()

}

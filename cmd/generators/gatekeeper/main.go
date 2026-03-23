package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"

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

// type Inventory struct {
// 	Material string
// 	Count    uint
// }
// sweaters := Inventory{"wool", 17}
// tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
// if err != nil { panic(err) }
// err = tmpl.Execute(os.Stdout, sweaters)
// if err != nil { panic(err) }

func main() {

	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	path := os.Args[1]

	fmt.Printf("path = %s\n", path)

	yamlBytes, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	t := Config{}

	yaml.Unmarshal(yamlBytes, &t)

	tmpl, err := template.New("structBuilder").Parse(
		`
package {{ .Package }}

import (
	"net/http"
	"strconv"
	"fmt"
)

{{ range .Objects }}
type {{ .Name }} struct {
{{ range .Fields }}	{{ .Name }} {{ .Kind }}
{{ end }}	validators []Validator
}

func New{{ .Name }}(r *http.Request) (*{{ .Name }}, []error) {
	errs := []error{}
	var err error
	x := new({{.Name }})
{{ range .Fields }}
{{ if eq .Kind "string" }}	x.{{ .Name }} = r.FormValue("{{ .Name }}"){{ end }}{{ if eq .Kind "int" }}	x.{{ .Name }}, err = strconv.Atoi(r.FormValue("{{ .Name }}"))
        if err != nil {
		errs = append(errs, fmt.Errorf("{{ .Name }} must be a number"))
        }
      {{ end }}
      {{ end }}
	
      return x, errs

}

{{ end }}



`)

	// func (s *{{ .Name }}) FromRequest(r *http.Request) {
	//   s.{{ .Name }} = r.FormValue("{{ .Name }}")
	//
	// {{ end }}

	if err != nil {
		panic(err)
	}

	fmt.Printf("t.Objects[0] = %v\n", t.Objects[0].Fields[0])

	file, err := os.OpenFile("gatekeeper.go", os.O_WRONLY|os.O_CREATE, 0777)

	err = tmpl.Execute(file, t)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// os.WriteFile()

	// fmt.Fprint()

}

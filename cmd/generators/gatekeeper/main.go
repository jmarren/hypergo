package main

import (
	"flag"
	"fmt"
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

type CustomValidator struct {
	Package  string `yaml:"package"`
	Function string `yaml:"function"`
}

type Validators struct {
	MinLen  string            `yaml:"minLen"`
	MaxLen  string            `yaml:"maxLen"`
	Min     string            `yaml:"min"`
	Max     string            `yaml:"max"`
	Options []string          `yaml:"options"`
	Email   bool              `yaml:"email"`
	Custom  []CustomValidator `yaml:"custom"`
}

type Field struct {
	Name       string
	Kind       string
	FormName   string `yaml:"formName"`
	Validators `yaml:"validators,inline"`
}

type Object struct {
	Name   string
	Fields []Field `yaml:"fields"`
}

type Config struct {
	Package string
	Objects []Object `yaml:"objects"`
}

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

func joinStrings(strs []string) string {

	quotedStrs := []string{}
	for _, str := range strs {
		quotedStrs = append(quotedStrs, "\""+str+"\"")
	}

	return strings.Join(quotedStrs, ", ")
}

func readYaml() Config {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	path := os.Args[1]

	yamlBytes, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	t := Config{}

	err = yaml.Unmarshal(yamlBytes, &t)

	if err != nil {
		panic(err)
	}
	return t

}

func main() {

	t := readYaml()

	fmt.Printf("t = %v\n", t)

	fmt.Printf("t.Objects[0].Fields = %v\n", t.Objects[0].Fields)

	templates := buildTemplates()

	file, err := os.OpenFile("gatekeeper.go", os.O_WRONLY|os.O_CREATE, 0777)

	err = templates.ExecuteTemplate(file, "base", t)
	if err != nil {
		panic(err)
	}
	defer file.Close()

}

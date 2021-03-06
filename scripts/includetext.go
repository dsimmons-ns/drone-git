// +build ignore

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	input  stringSlice
	output string
	name   string
)

func main() {
	flag.Var(&input, "input", "input files")
	flag.StringVar(&output, "output", "", "output file")
	flag.StringVar(&name, "package", "", "package name")
	flag.Parse()

	var files []File
	for _, file := range input {
		out, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		files = append(files, File{
			Name: file,
			Slug: slugify(file),
			Data: string(out),
		})
	}

	data := map[string]interface{}{
		"Files":   files,
		"Package": name,
	}
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
	if err != nil {
		log.Fatalln(err)
	}

	ioutil.WriteFile(output, buf.Bytes(), 0644)
}

func slugify(s string) string {
	ext := filepath.Ext(s)
	s = strings.TrimSuffix(s, ext)
	s = strings.Title(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type File struct {
	Name string
	Data string
	Slug string
}

var tmpl = template.Must(template.New("_").Parse(`package {{ .Package }}

// DO NOT EDIT. This file is automatically generated.

{{ range .Files -}}
// Contents of {{ .Name }}
const {{ .Slug }} = ` + "`{{ .Data }}`" + `

{{ end -}}`))

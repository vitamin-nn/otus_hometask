package main

import (
	"bytes"
	"go/format"
	"os"
	"strings"
	"text/template"
)

func generate(outFilePath string, parsedData ParsedData) error {
	f, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	tmpl := template.Must(
		template.New("main").
			Funcs(template.FuncMap{"getFieldVar": strings.ToLower}).
			Parse(getTemplates()),
	)
	var b bytes.Buffer
	err = tmpl.ExecuteTemplate(&b, "validation", parsedData)
	if err != nil {
		return err
	}
	bGen, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}
	_, err = f.Write(bGen)
	return err
}

package main

import (
	"os"
	"text/template"
)

type TestData struct {
	Type   string
	Server []string
}

func template_engine(data TestData) {
	tmpl, err := template.ParseFiles("test.tmlp")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

type context struct {
}

func (c *context) Env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}

func createConfig(templateFile string, outputFile string) {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(file, &context{})
	file.Close()
}

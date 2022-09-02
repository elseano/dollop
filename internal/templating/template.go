package templating

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

func BuildTemplateText(value string) *template.Template {
	tmpl, err := template.New("tmpl").Funcs(templateFuncs).Parse(value)
	if err != nil {
		log.Fatalf("Cannot build template for \"%s\": %s", value, err.Error())
		panic("Invalid config")
	}

	return tmpl.Option("missingkey=invalid")
}

func BuildTemplate(value string) *template.Template {
	if strings.Contains(value, "{{") {
		return BuildTemplateText(value)
	}

	return BuildTemplateText("{{ ." + value + " }}")
}

func ApplyTemplate(tmpl *template.Template, data interface{}) (string, error) {
	b := bytes.Buffer{}
	err := tmpl.Execute(&b, data)

	s := b.String()

	if strings.Contains(s, "<no value>") {
		s = ""
	}

	return s, err
}

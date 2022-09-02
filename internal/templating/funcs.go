package templating

import (
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var templateFuncs = template.FuncMap{
	"Contains":      strings.Contains,
	"ContainsAny":   strings.ContainsAny,
	"Fields":        strings.Fields,
	"HasPrefix":     strings.HasPrefix,
	"HasSuffix":     strings.HasSuffix,
	"StrIndex":      strings.Index,
	"StrJoin":       strings.Join,
	"Replace":       strings.Replace,
	"ReplaceAll":    strings.ReplaceAll,
	"Split":         strings.Split,
	"SplitN":        strings.SplitN,
	"Title":         titleCase,
	"ToLower":       strings.ToLower,
	"ToUpper":       strings.ToUpper,
	"TrimPrefix":    strings.TrimPrefix,
	"TrimSpace":     strings.TrimSpace,
	"TrimSuffix":    strings.TrimSuffix,
	"TruncateLeft":  truncateLeft,
	"Truncate":      truncateRight,
	"FormatSeconds": formatSeconds,
	"div":           div,
	"mul":           mul,
	"add":           add,
	"sub":           sub,
}

func titleCase(str string) string {
	return cases.Title(language.English).String(str)
}

func truncateLeft(str string, limit int) string {
	if len(str) > limit {
		return "..." + str[len(str)-limit:]
	}

	return str
}

func truncateRight(str string, limit int) string {
	if len(str) > limit {
		return str[0:limit] + "..."
	}

	return str
}

func formatSeconds(seconds float64) string {
	var duration = time.Duration(seconds * 1000000000)
	return duration.String()
}

func mul(a float64, b float64) float64 {
	return a * b
}

func div(a float64, b float64) float64 {
	return a / b
}

func add(a float64, b float64) float64 {
	return a + b
}

func sub(a float64, b float64) float64 {
	return a - b
}

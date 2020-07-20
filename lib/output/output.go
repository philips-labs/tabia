package output

import (
	"encoding/json"
	"io"
	"text/template"
)

// PrintJSON prints the json using indentation using the given writer
func PrintJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	return enc.Encode(data)
}

// PrintUsingTemplate prints the data using the given template
func PrintUsingTemplate(w io.Writer, templateContent []byte, data interface{}) error {
	tmpl, err := template.New("template").Parse(string(templateContent))
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

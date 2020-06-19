package output

import (
	"encoding/json"
	"io"
)

// PrintJSON prints the json using indentation using the given writer
func PrintJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	err := enc.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

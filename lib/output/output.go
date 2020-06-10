package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// PrintJSON prints the json using indentation using the given writer
func PrintJSON(w io.Writer, data interface{}) error {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", json)
	return nil
}

package output_test

import (
	"strings"
	"testing"

	"github.com/philips-labs/tabia/lib/output"

	"github.com/stretchr/testify/assert"
)

func TestPrintJson(t *testing.T) {
	assert := assert.New(t)

	data := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    int    `json:"prio"`
	}{Title: "JSON printing", Description: "Prints JSON to io.Writer", Priority: 1}

	var builder strings.Builder
	err := output.PrintJSON(&builder, data)

	if assert.NoError(err) {
		assert.Equal("{\n  \"title\": \"JSON printing\",\n  \"description\": \"Prints JSON to io.Writer\",\n  \"prio\": 1\n}\n", builder.String())
	}
}

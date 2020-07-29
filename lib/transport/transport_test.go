package transport_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/transport"
)

func TestTeeRoundTripper(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	}))
	defer ts.Close()

	var writer strings.Builder
	client := http.Client{
		Transport: transport.TeeRoundTripper{
			RoundTripper: http.DefaultTransport,
			Writer:       &writer,
		},
	}

	jsonString := `{ "say": "hello-world", "to": "marco" }`
	json := strings.NewReader(jsonString)
	_, err := client.Post(ts.URL, "application/json", json)

	assert.NoError(err)
	assert.NotEmpty(writer.String())
	assert.Equal(fmt.Sprintf("POST: %s %s", ts.URL, jsonString), writer.String())

	writer.Reset()
	_, err = client.Get(ts.URL)

	assert.NoError(err)
	assert.NotEmpty(writer.String())
	assert.Equal(fmt.Sprintf("GET: %s ", ts.URL), writer.String())

}

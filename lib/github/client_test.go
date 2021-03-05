package github_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder

	client := github.NewClientWithTokenAuth(os.Getenv("TABIA_GITHUB_TOKEN"), &buf)
	var q struct{}
	err := client.Client.Query(context.Background(), q, nil)

	assert.EqualError(err, "Field must have selections (anonymous query returns Query but has no selections. Did you mean ' { ... }'?)")
	assert.NotEmpty(buf)
	assert.Equal("POST: https://api.github.com/graphql {\"query\":\"{}\"}\n", buf.String())
}

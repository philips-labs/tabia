package github_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder

	client := github.NewClientWithTokenAuth("token", &buf)
	var q struct{}
	client.Client.Query(context.Background(), q, nil)

	assert.NotEmpty(buf)
	assert.Equal("POST: https://api.github.com/graphql {\"query\":\"{}\"}\n", buf.String())
}

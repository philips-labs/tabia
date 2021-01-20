package github_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
)

func TestAppClient(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder
	integrationID := int64(12345)
	client, err := github.NewClientWithAppAuth(integrationID, "/path/to/rsa-private-key.pem", &buf)
	assert.Error(err)
	assert.Nil(client)
}

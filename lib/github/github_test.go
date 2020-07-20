package github_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
)

func TestRepositoryVisibility(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Public", github.Public.String())
	assert.Equal("Internal", github.Internal.String())
	assert.Equal("Private", github.Private.String())
}

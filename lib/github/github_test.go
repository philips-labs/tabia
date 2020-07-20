package github_test

import (
	"encoding/json"
	"fmt"
	"strings"
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

func TestRepositoryVisibilityToJSON(t *testing.T) {
	assert := assert.New(t)

	expectedTemplate := `{"name":"%s","visibility":"%s","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","pushed_at":"0001-01-01T00:00:00Z"}
`

	var result strings.Builder
	jsonEnc := json.NewEncoder(&result)

	privRepo := github.Repository{
		Name:       "private-repo",
		Visibility: github.Private,
	}
	jsonEnc.Encode(privRepo)
	assert.Equal(fmt.Sprintf(expectedTemplate, "private-repo", "Private"), result.String())

	internalRepo := github.Repository{
		Name:       "internal-repo",
		Visibility: github.Internal,
	}
	result.Reset()
	jsonEnc.Encode(internalRepo)
	assert.Equal(fmt.Sprintf(expectedTemplate, "internal-repo", "Internal"), result.String())

	publicRepo := github.Repository{
		Name:       "public-repo",
		Visibility: github.Public,
	}
	result.Reset()
	jsonEnc.Encode(publicRepo)
	assert.Equal(fmt.Sprintf(expectedTemplate, "public-repo", "Public"), result.String())
}

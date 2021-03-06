package shared_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/shared"
)

func TestRepositoryVisibility(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Public", shared.Public.String())
	assert.Equal("Internal", shared.Internal.String())
	assert.Equal("Private", shared.Private.String())
}

func TestVisibilityFromText(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(shared.Public, shared.VisibilityFromText(""))
	assert.Equal(shared.Public, shared.VisibilityFromText("unknown"))
	assert.Equal(shared.Public, shared.VisibilityFromText("public"))
	assert.Equal(shared.Public, shared.VisibilityFromText("PUBLIC"))
	assert.Equal(shared.Public, shared.VisibilityFromText("Public"))
	assert.Equal(shared.Internal, shared.VisibilityFromText("internal"))
	assert.Equal(shared.Internal, shared.VisibilityFromText("INTERNAL"))
	assert.Equal(shared.Internal, shared.VisibilityFromText("Internal"))
	assert.Equal(shared.Private, shared.VisibilityFromText("private"))
	assert.Equal(shared.Private, shared.VisibilityFromText("PRIVATE"))
	assert.Equal(shared.Private, shared.VisibilityFromText("Private"))
}

func TestVisibilityMarshalling(t *testing.T) {
	assert := assert.New(t)
	blob := `["Public","Internal","Private"]`

	var visibilities []shared.Visibility

	err := json.Unmarshal([]byte(blob), &visibilities)
	if assert.NoError(err) {
		newBlob, err := json.Marshal(visibilities)
		if assert.NoError(err) {
			assert.Equal(blob, string(newBlob))
		}
	}
}

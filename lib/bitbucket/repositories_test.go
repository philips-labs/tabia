package bitbucket_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-labs/tabia/lib/bitbucket"

	"github.com/stretchr/testify/assert"
)

var stubRepositoriesResponse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	projectsResponse := bitbucket.RepositoriesResponse{
		Values: make([]bitbucket.Repository, 12),
	}
	resp, _ := json.Marshal(projectsResponse)
	w.Write(resp)
})

func TestListRepositoriesRaw(t *testing.T) {
	assert := assert.New(t)

	bb, apiBaseURL, teardown := bitbucketTestClient(stubRepositoriesResponse)
	defer teardown()
	resp, err := bb.RawRequest("GET", apiBaseURL+"/projects/VID/repos", "")
	if !assert.NoError(err) {
		return
	}
	defer resp.Close()

	assert.NotNil(resp)
	bytes, err := io.ReadAll(resp)

	assert.NotEmpty(bytes)
}

func TestListRepositories(t *testing.T) {
	assert := assert.New(t)

	bb, _, teardown := bitbucketTestClient(stubRepositoriesResponse)
	defer teardown()
	resp, err := bb.Repositories.List("ACE")

	if !assert.NoError(err) {
		return
	}

	assert.NotNil(resp)
	assert.Len(resp.Values, 12)
}

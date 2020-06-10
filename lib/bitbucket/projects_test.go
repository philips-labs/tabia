package bitbucket_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/philips-labs/tabia/lib/bitbucket"

	"github.com/stretchr/testify/assert"
)

var stubProjectsResponse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	projectsResponse := bitbucket.ProjectsResponse{
		Values: make([]bitbucket.Project, 100),
	}
	resp, _ := json.Marshal(projectsResponse)
	w.Write(resp)
})

func TestListProjectsRaw(t *testing.T) {
	assert := assert.New(t)
	bb, apiBaseURL, teardown := bitbucketTestClient(stubProjectsResponse)
	defer teardown()

	resp, err := bb.RawRequest("GET", apiBaseURL+"/projects", "")
	if !assert.NoError(err) {
		return
	}
	defer resp.Close()

	assert.NotNil(resp)
	bytes, err := ioutil.ReadAll(resp)

	assert.NotEmpty(bytes)
}

func TestListProjects(t *testing.T) {
	assert := assert.New(t)

	bb, _, teardown := bitbucketTestClient(stubProjectsResponse)
	defer teardown()

	resp, err := bb.Projects.List(0)

	if !assert.NoError(err) {
		return
	}

	assert.NotNil(resp)
	assert.Len(resp.Values, 100)
}

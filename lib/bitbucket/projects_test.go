package bitbucket_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListProjectsRaw(t *testing.T) {
	assert := assert.New(t)

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

	resp, err := bb.Projects.List(0)

	if !assert.NoError(err) {
		return
	}

	assert.NotNil(resp)
	assert.Len(resp.Values, 100)
}

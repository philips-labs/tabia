package bitbucket_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRepositoriesRaw(t *testing.T) {
	assert := assert.New(t)

	resp, err := bb.RawRequest("GET", apiBaseURL+"/projects/VID/repos", "")
	if !assert.NoError(err) {
		return
	}
	defer resp.Close()

	assert.NotNil(resp)
	bytes, err := ioutil.ReadAll(resp)

	assert.NotEmpty(bytes)
}

func TestListRepositories(t *testing.T) {
	assert := assert.New(t)

	resp, err := bb.Repositories.List("ACE")

	if !assert.NoError(err) {
		return
	}

	assert.NotNil(resp)
	assert.Len(resp.Values, 1)
}

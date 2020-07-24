package github_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
)

func TestDownloadContents(t *testing.T) {
	assert := assert.New(t)
	gh := github.NewClientWithTokenAuth(os.Getenv("TABIA_GITHUB_TOKEN"))
	contents, err := gh.DownloadContents(context.Background(), "philips-labs", "tabia", "README.md")
	if assert.NoError(err) {
		readme, _ := ioutil.ReadFile("../../README.md")
		assert.NotEmpty(contents)
		assert.Equal(string(readme[:100]), string(contents[:100]))
	}

	contents, err = gh.DownloadContents(context.Background(), "philips-labs", "tabia", "IamNotThere.txt")
	if assert.Error(err) {
		assert.EqualError(err, "No file named IamNotThere.txt found in .")
	}
}

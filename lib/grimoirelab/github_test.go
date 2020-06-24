package grimoirelab_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
	"github.com/philips-labs/tabia/lib/grimoirelab"
)

func TestConvertGithubProjectsJSON(t *testing.T) {
	type MyString string

	assert := assert.New(t)

	ghUser := os.Getenv("TABIA_GITHUB_USER")
	ghToken := os.Getenv("TABIA_GITHUB_TOKEN")

	repos := []github.Repository{
		github.Repository{
			Name:      "R1",
			IsPrivate: false,
			URL:       "https://github.com/philips-software/logproxy",
			Owner: github.Owner{
				Login: "philips-software",
			},
		},
		github.Repository{
			Name:      "R1",
			IsPrivate: true,
			URL:       "https://github.com/philips-labs/tabia",
			Owner: github.Owner{
				Login: "philips-labs",
			},
		},
	}

	projects := grimoirelab.ConvertGithubToProjectsJSON(
		repos, func(repo github.Repository) grimoirelab.Metadata {
			return grimoirelab.Metadata{
				"title":   repo.Owner.Login,
				"program": "One Codebase",
			}
		})

	if assert.Len(projects, 2) {
		if assert.Len(projects["philips-software"].Git, 1) {
			assert.Equal("https://github.com/philips-software/logproxy.git", projects["philips-software"].Git[0])
			assert.Equal("https://github.com/philips-software/logproxy", projects["philips-software"].Github[0])
			assert.Equal("https://github.com/philips-software/logproxy", projects["philips-software"].GithubRepo[0])
		}
		assert.Len(projects["philips-software"].Metadata, 2)

		if assert.Len(projects["philips-labs"].Git, 1) {
			assertUrlHasBasicAuth(t, projects["philips-labs"].Git[0], "https", ghUser, ghToken, "github.com", "/philips-labs/tabia.git")
		}
		if assert.Len(projects["philips-labs"].Github, 1) {
			assertUrlHasBasicAuth(t, projects["philips-labs"].Github[0], "https", ghUser, ghToken, "github.com", "/philips-labs/tabia")
		}
		if assert.Len(projects["philips-labs"].GithubRepo, 1) {
			assertUrlHasBasicAuth(t, projects["philips-labs"].GithubRepo[0], "https", ghUser, ghToken, "github.com", "/philips-labs/tabia")
		}
		assert.Len(projects["philips-labs"].Metadata, 2)
	}
}

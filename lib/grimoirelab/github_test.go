package grimoirelab_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
	"github.com/philips-labs/tabia/lib/grimoirelab"
	"github.com/philips-labs/tabia/lib/shared"
)

func TestNewGithubProjectMatcherFromJSON(t *testing.T) {
	assert := assert.New(t)

	json := strings.NewReader(`{
	"rules": {
		"My Project": { "url": "(?i)foo|Bar|BAZ" }
	}
}`)
	m, err := grimoirelab.NewGithubProjectMatcherFromJSON(json)
	if assert.NoError(err) {
		assert.Equal("(?i)foo|Bar|BAZ", m.Rules["My Project"].URL.String())
	}

	json = strings.NewReader(`{
		"rules": {
			"My Project": { "url": "" },
		}
	}`)
	m, err = grimoirelab.NewGithubProjectMatcherFromJSON(json)
	assert.EqualError(err, "invalid character '}' looking for beginning of object key string")
	assert.Nil(m)

	json = strings.NewReader(`{
		"rules": {
			"My Project": { "url": "(invalid|regex" }
		}
	}`)
	m, err = grimoirelab.NewGithubProjectMatcherFromJSON(json)
	assert.EqualError(err, "error parsing regexp: missing closing ): `(invalid|regex`")
	assert.Nil(m)
}

func TestConvertGithubProjectsJSON(t *testing.T) {
	assert := assert.New(t)

	ghUser := os.Getenv("TABIA_GITHUB_USER")
	ghToken := os.Getenv("TABIA_GITHUB_TOKEN")

	repos := []github.Repository{
		{
			Name:       "R1",
			Visibility: shared.Public,
			URL:        "https://github.com/philips-software/logproxy",
			Owner:      "philips-software",
		},
		{
			Name:       "R1",
			Visibility: shared.Private,
			URL:        "https://github.com/philips-labs/tabia",
			Owner:      "philips-labs",
		},
	}

	projects := grimoirelab.ConvertGithubToProjectsJSON(
		repos,
		func(repo github.Repository) grimoirelab.Metadata {
			return grimoirelab.Metadata{
				"title":   repo.Owner,
				"program": "One Codebase",
			}
		},
		&grimoirelab.GithubProjectMatcher{
			Rules: map[string]grimoirelab.GithubProjectMatcherRule{
				"One Codebase": {URL: grimoirelab.MustCompile("(?i)Tabia")},
			},
		},
	)

	if assert.Len(projects, 2) {
		if assert.Len(projects["philips-software"].Git, 1) {
			assert.Equal("https://github.com/philips-software/logproxy.git", projects["philips-software"].Git[0])
			assert.Equal("https://github.com/philips-software/logproxy", projects["philips-software"].Github[0])
			assert.Equal("https://github.com/philips-software/logproxy", projects["philips-software"].GithubRepo[0])
		}
		assert.Len(projects["philips-software"].Metadata, 2)

		if assert.Len(projects["One Codebase"].Git, 1) {
			assertUrlHasBasicAuth(t, projects["One Codebase"].Git[0], "https", ghUser, ghToken, "github.com", "/philips-labs/tabia.git")
		}
		if assert.Len(projects["One Codebase"].Github, 1) {
			assertUrlHasBasicAuth(t, projects["One Codebase"].Github[0], "https", ghUser, ghToken, "github.com", "/philips-labs/tabia")
		}
		if assert.Len(projects["One Codebase"].GithubRepo, 1) {
			assertUrlHasBasicAuth(t, projects["One Codebase"].GithubRepo[0], "https", ghUser, ghToken, "github.com", "/philips-labs/tabia")
		}
		assert.Len(projects["One Codebase"].Metadata, 2)
	}
}

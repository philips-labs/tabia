package grimoirelab_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/bitbucket"
	"github.com/philips-labs/tabia/lib/github"
	"github.com/philips-labs/tabia/lib/grimoirelab"
)

func TestConvertBitbucketProjectsJSON(t *testing.T) {
	assert := assert.New(t)

	bbUser := os.Getenv("TABIA_BITBUCKET_USER")
	bbToken := os.Getenv("TABIA_BITBUCKET_TOKEN")

	repos := []bitbucket.Repository{
		bitbucket.Repository{
			Project: bitbucket.Project{Name: "P1"},
			Name:    "R1",
			Public:  true,
			Links: bitbucket.Links{
				Clone: []bitbucket.CloneLink{
					bitbucket.CloneLink{Name: "http", Href: "https://bitbucket.org/scm/p1/r1.git"},
				},
			},
		},
		bitbucket.Repository{
			Project: bitbucket.Project{Name: "P1"},
			Name:    "R2",
			Public:  false,
			Links: bitbucket.Links{
				Clone: []bitbucket.CloneLink{
					bitbucket.CloneLink{Name: "http", Href: "https://bitbucket.org/scm/p1/r2.git"},
				},
			},
		},
		bitbucket.Repository{
			Project: bitbucket.Project{Name: "P2"},
			Name:    "R1",
			Public:  false,
			Links: bitbucket.Links{
				Clone: []bitbucket.CloneLink{
					bitbucket.CloneLink{Name: "http", Href: "https://bitbucket.org/scm/p2/r1.git"},
				},
			},
		},
		bitbucket.Repository{
			Project: bitbucket.Project{Name: "P2"},
			Name:    "R2",
			Public:  false,
			Links: bitbucket.Links{
				Clone: []bitbucket.CloneLink{
					bitbucket.CloneLink{Name: "http", Href: "https://bitbucket.org/scm/p2/r2.git"},
				},
			},
		},
		bitbucket.Repository{
			Project: bitbucket.Project{Name: "P2"},
			Name:    "R3",
			Public:  true,
			Links: bitbucket.Links{
				Clone: []bitbucket.CloneLink{
					bitbucket.CloneLink{Name: "http", Href: "https://bitbucket.org/scm/p2/r3.git"},
				},
			},
		},
	}

	projects := grimoirelab.ConvertBitbucketToProjectsJSON(repos, func(repo bitbucket.Repository) grimoirelab.Metadata {
		return grimoirelab.Metadata{
			"title":   repo.Project.Name,
			"program": "One Codebase",
		}
	})

	if assert.Len(projects, 2) {
		gitP1 := projects["P1"].Git
		if assert.Len(gitP1, 2) {
			assert.Equal(gitP1[0], "https://bitbucket.org/scm/p1/r1.git")
			assertUrlHasBasicAuth(t, gitP1[1], "https", bbUser, bbToken, "bitbucket.org", "/scm/p1/r2.git")
		}
		assert.Len(projects["P1"].Metadata, 2)

		gitP2 := projects["P2"].Git
		if assert.Len(gitP2, 3) {
			assertUrlHasBasicAuth(t, gitP2[0], "https", bbUser, bbToken, "bitbucket.org", "/scm/p2/r1.git")
			assertUrlHasBasicAuth(t, gitP2[1], "https", bbUser, bbToken, "bitbucket.org", "/scm/p2/r2.git")
			assert.Equal(gitP2[2], "https://bitbucket.org/scm/p2/r3.git")
		}
		assert.Len(projects["P2"].Metadata, 2)
	}
}

func TestConvertGithubProjectsJSON(t *testing.T) {
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

	projects := grimoirelab.ConvertGithubToProjectsJSON(repos, func(repo github.Repository) grimoirelab.Metadata {
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

func assertUrlHasBasicAuth(t *testing.T, uri, scheme, user, password, hostname, path string) {
	assert := assert.New(t)
	u, err := url.Parse(uri)
	assert.NoError(err)
	assert.Equal(scheme, u.Scheme)
	assert.Equal(user, u.User.Username())
	pass, isSet := u.User.Password()
	assert.True(isSet)
	assert.Equal(password, pass)
	assert.Equal(hostname, u.Hostname())
	assert.Equal(path, u.EscapedPath())
}

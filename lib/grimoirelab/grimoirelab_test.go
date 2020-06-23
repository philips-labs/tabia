package grimoirelab_test

import (
	"fmt"
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
	basicAuth := fmt.Sprintf("%s:%s", bbUser, bbToken)

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
			assert.Equal("https://", gitP1[1][:8])
			assert.Contains(gitP1[1], basicAuth)
			assert.Contains(gitP1[1], "@bitbucket.org/scm/p1/r2.git")
		}
		assert.Len(projects["P1"].Metadata, 2)

		gitP2 := projects["P2"].Git
		if assert.Len(gitP2, 3) {
			assert.Equal("https://", gitP2[0][:8])
			assert.Contains(gitP2[0], basicAuth)
			assert.Contains(gitP2[0], "@bitbucket.org/scm/p2/r1.git")
			assert.Equal("https://", gitP2[1][:8])
			assert.Contains(gitP2[1], basicAuth)
			assert.Contains(gitP2[1], "@bitbucket.org/scm/p2/r2.git")
			assert.Equal(gitP2[2], "https://bitbucket.org/scm/p2/r3.git")
		}
		assert.Len(projects["P2"].Metadata, 2)
	}
}

func TestConvertGithubProjectsJSON(t *testing.T) {
	assert := assert.New(t)

	ghUser := os.Getenv("TABIA_GITHUB_USER")
	ghToken := os.Getenv("TABIA_GITHUB_TOKEN")
	basicAuth := fmt.Sprintf("%s:%s", ghUser, ghToken)

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
		}
		assert.Len(projects["philips-software"].Metadata, 2)

		if assert.Len(projects["philips-labs"].Git, 1) {
			assert.Equal("https://", projects["philips-labs"].Git[0][:8])
			assert.Contains(projects["philips-labs"].Git[0], basicAuth)
			assert.Contains(projects["philips-labs"].Git[0], "@github.com/philips-labs/tabia.git")
			assert.Equal("https://", projects["philips-labs"].Github[0][:8])
			assert.Contains(projects["philips-labs"].Github[0], basicAuth)
			assert.Contains(projects["philips-labs"].Github[0], "@github.com/philips-labs/tabia")
		}
		assert.Len(projects["philips-labs"].Metadata, 2)
	}
}

package grimoirelab_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/bitbucket"
	"github.com/philips-labs/tabia/lib/grimoirelab"
)

func TestConvertProjectsJSON(t *testing.T) {
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

	projects := grimoirelab.ConvertProjectsJSON(repos, func(repo bitbucket.Repository) grimoirelab.Metadata {
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

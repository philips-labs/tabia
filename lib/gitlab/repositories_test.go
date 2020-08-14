package gitlab_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gl "github.com/xanzy/go-gitlab"

	"github.com/philips-labs/tabia/lib/gitlab"
)

func TestMap(t *testing.T) {
	assert := assert.New(t)

	projects := []*gl.Project{
		&gl.Project{},
	}

	repos := gitlab.Map(projects)

	assert.Len(repos, len(projects))
}

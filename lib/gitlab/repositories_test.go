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

func TestBuildListProjectsOptions(t *testing.T) {
	assert := assert.New(t)

	opt := gitlab.BuildListProjectsOptions()
	assert.Equal(gl.ListOptions{
		PerPage: 100,
		Page:    1,
	}, opt.ListOptions)
	assert.Equal((*gl.VisibilityValue)(nil), opt.Visibility)

	opt = gitlab.BuildListProjectsOptions(gitlab.WithPublicVisibility)
	assert.Equal(gl.ListOptions{
		PerPage: 100,
		Page:    1,
	}, opt.ListOptions)
	assert.Equal(visibilityValuePtr(gl.PublicVisibility), opt.Visibility)

	opt = gitlab.BuildListProjectsOptions(gitlab.WithInternalVisibility)
	assert.Equal(gl.ListOptions{
		PerPage: 100,
		Page:    1,
	}, opt.ListOptions)
	assert.Equal(visibilityValuePtr(gl.InternalVisibility), opt.Visibility)

	opt = gitlab.BuildListProjectsOptions(gitlab.WithPrivateVisibility)
	assert.Equal(gl.ListOptions{
		PerPage: 100,
		Page:    1,
	}, opt.ListOptions)
	assert.Equal(visibilityValuePtr(gl.PrivateVisibility), opt.Visibility)
}

func visibilityValuePtr(v gl.VisibilityValue) *gl.VisibilityValue {
	return &v
}

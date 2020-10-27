package gitlab_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/gitlab"
	"github.com/philips-labs/tabia/lib/shared"
)

func TestReduce(t *testing.T) {
	assert := assert.New(t)

	repos := []gitlab.Repository{
		gitlab.Repository{
			ID:   1,
			Name: "cool", Visibility: shared.Public, Owner: "research",
			CreatedAt:      timePointer(time.Now().Add(-24 * time.Hour)),
			LastActivityAt: timePointer(time.Now().Add(-24 * time.Hour)),
		},
		gitlab.Repository{
			ID:   2,
			Name: "stuff", Visibility: shared.Public, Owner: "research",
			CreatedAt:      timePointer(time.Now().Add(-96 * time.Hour)),
			LastActivityAt: timePointer(time.Now().Add(-96 * time.Hour)),
		},
		gitlab.Repository{
			ID:   3,
			Name: "happens", Visibility: shared.Public, Owner: "research",
			CreatedAt:      timePointer(time.Now().Add(-24 * time.Hour)),
			LastActivityAt: timePointer(time.Now().Add(-24 * time.Hour)),
		},
		gitlab.Repository{
			ID:   4,
			Name: "at", Visibility: shared.Internal, Owner: "research",
			CreatedAt:      timePointer(time.Now().Add(-48 * time.Hour)),
			LastActivityAt: timePointer(time.Now().Add(-48 * time.Hour)),
		},
		gitlab.Repository{
			ID:   5,
			Name: "philips", Visibility: shared.Private, Owner: "research",
			CreatedAt:      timePointer(time.Now().Add(-24 * time.Hour)),
			LastActivityAt: timePointer(time.Now().Add(-24 * time.Hour)),
		},
	}

	reduced, err := gitlab.Reduce(repos, "")
	if assert.NoError(err) {
		assert.Len(reduced, 5)
		assert.ElementsMatch(reduced, repos)
	}

	reduced, err = gitlab.Reduce(repos, `{ .Name == "stuff" }`)
	if assert.NoError(err) {
		assert.Len(reduced, 1)
		assert.Contains(reduced, repos[1])
	}

	reduced, err = gitlab.Reduce(repos, `{ .IsPublic() }`)
	if assert.NoError(err) {
		assert.Len(reduced, 3)
		assert.Contains(reduced, repos[0])
		assert.Contains(reduced, repos[1])
		assert.Contains(reduced, repos[2])
	}

	reduced, err = gitlab.Reduce(repos, `{ .IsInternal() }`)
	if assert.NoError(err) {
		assert.Len(reduced, 1)
		assert.Contains(reduced, repos[3])
	}

	reduced, err = gitlab.Reduce(repos, `{ .IsPrivate() }`)
	if assert.NoError(err) {
		assert.Len(reduced, 1)
		assert.Contains(reduced, repos[4])
	}

	reduced, err = gitlab.Reduce(repos, `{ Contains(.Name, "a") }`)
	if assert.NoError(err) {
		assert.Len(reduced, 2)
		assert.Contains(reduced, repos[2])
		assert.Contains(reduced, repos[3])
	}

	since := time.Now().Add(-25 * time.Hour).Format(time.RFC3339)
	reduced, err = gitlab.Reduce(repos, fmt.Sprintf(`{ .CreatedSince("%s") }`, since))
	if assert.NoError(err) {
		assert.Len(reduced, 3)
		assert.Contains(reduced, repos[0])
		assert.Contains(reduced, repos[2])
		assert.Contains(reduced, repos[4])
	}

	since = time.Now().Add(-97 * time.Hour).Format(time.RFC3339)
	reduced, err = gitlab.Reduce(repos, fmt.Sprintf(`{ .LastActivitySince("%s") }`, since))
	if assert.NoError(err) {
		assert.Len(reduced, 5)
		assert.Contains(reduced, repos[0])
		assert.Contains(reduced, repos[1])
		assert.Contains(reduced, repos[2])
		assert.Contains(reduced, repos[3])
		assert.Contains(reduced, repos[4])
	}
}

func TestReduceWrongExpression(t *testing.T) {
	assert := assert.New(t)

	repos := []gitlab.Repository{
		gitlab.Repository{Name: "cool", Visibility: shared.Public},
	}

	reduced, err := gitlab.Reduce(repos, `.Name == "cool"`)
	assert.Error(err)
	assert.EqualError(err, "unexpected token Operator(\".\") (1:22)\n | filter(Repositories, .Name == \"cool\")\n | .....................^")
	assert.Nil(reduced)

	reduced, err = gitlab.Reduce(repos, `{ UnExistingFunc(.URL, "stuff") }`)
	assert.Error(err)
	assert.EqualError(err, "cannot get \"UnExistingFunc\" from gitlab.RepositoryFilterEnv (1:24)\n | filter(Repositories, { UnExistingFunc(.URL, \"stuff\") })\n | .......................^")
	assert.Nil(reduced)
}

func timePointer(t time.Time) *time.Time {
	return &t
}

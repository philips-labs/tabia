package github_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
)

func TestReduce(t *testing.T) {
	assert := assert.New(t)

	repos := []github.Repository{
		github.Repository{
			Name: "tabia", Visibility: github.Public, Owner: "philips-labs",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			PushedAt:  time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		github.Repository{
			Name: "garo", Visibility: github.Public, Owner: "philips-labs",
			CreatedAt: time.Now().Add(-96 * time.Hour),
			PushedAt:  time.Now().Add(-96 * time.Hour),
			UpdatedAt: time.Now().Add(-96 * time.Hour),
		},
		github.Repository{
			Name: "dct-notary-admin", Visibility: github.Public, Owner: "philips-labs",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			PushedAt:  time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		github.Repository{
			Name: "company-draft", Visibility: github.Internal, Owner: "philips-labs",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			PushedAt:  time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-48 * time.Hour),
		},
		github.Repository{
			Name: "top-secret", Visibility: github.Private, Owner: "philips-labs",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			PushedAt:  time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
			Topics:    []github.Topic{github.Topic{Name: "ip"}},
		},
	}

	reduced, err := github.Reduce(repos, "")
	if assert.NoError(err) {
		assert.Len(reduced, 5)
		assert.ElementsMatch(reduced, repos)
	}

	reduced, err = github.Reduce(repos, `{ .Name == "garo" }`)
	if assert.NoError(err) {
		assert.Len(reduced, 1)
		assert.Contains(reduced, repos[1])
	}

	reduced, err = github.Reduce(repos, `{ .IsPublic() }`)
	if assert.NoError(err) {
		assert.Len(reduced, 3)
		assert.Contains(reduced, repos[0])
		assert.Contains(reduced, repos[1])
		assert.Contains(reduced, repos[2])
	}

	reduced, err = github.Reduce(repos, `{ .IsInternal() }`)
	if assert.NoError(err) {
		assert.Len(reduced, 1)
		assert.Contains(reduced, repos[3])
	}

	reduced, err = github.Reduce(repos, `{ .IsPrivate() }`)
	if assert.NoError(err) {
		assert.Len(reduced, 1)
		assert.Contains(reduced, repos[4])
	}

	reduced, err = github.Reduce(repos, `{ Contains(.Name, "ar") }`)
	if assert.NoError(err) {
		assert.Len(reduced, 2)
		assert.Contains(reduced, repos[1])
		assert.Contains(reduced, repos[2])
	}

	reduced, err = github.Reduce(repos, `{ .HasTopic("ip") }`)
	if assert.NoError(err) {
		assert.Len(reduced, 1)
		assert.Contains(reduced, repos[4])
	}

	since := time.Now().Add(-25 * time.Hour).Format(time.RFC3339)
	reduced, err = github.Reduce(repos, fmt.Sprintf(`{ .CreatedSince("%s") }`, since))
	if assert.NoError(err) {
		assert.Len(reduced, 3)
		assert.Contains(reduced, repos[0])
		assert.Contains(reduced, repos[2])
		assert.Contains(reduced, repos[4])
	}

	since = time.Now().Add(-97 * time.Hour).Format(time.RFC3339)
	reduced, err = github.Reduce(repos, fmt.Sprintf(`{ .UpdatedSince("%s") }`, since))
	if assert.NoError(err) {
		assert.Len(reduced, 5)
		assert.Contains(reduced, repos[0])
		assert.Contains(reduced, repos[1])
		assert.Contains(reduced, repos[2])
		assert.Contains(reduced, repos[3])
		assert.Contains(reduced, repos[4])
	}

	since = time.Now().Add(-49 * time.Hour).Format(time.RFC3339)
	reduced, err = github.Reduce(repos, fmt.Sprintf(`{ .PushedSince("%s") }`, since))
	if assert.NoError(err) {
		assert.Len(reduced, 4)
		assert.Contains(reduced, repos[0])
		assert.Contains(reduced, repos[2])
		assert.Contains(reduced, repos[3])
		assert.Contains(reduced, repos[4])
	}
}

func TestReduceWrongExpression(t *testing.T) {
	assert := assert.New(t)

	repos := []github.Repository{
		github.Repository{Name: "tabia", Visibility: github.Public},
	}

	reduced, err := github.Reduce(repos, `.Name == "tabia"`)
	assert.Error(err)
	assert.EqualError(err, "unexpected token Operator(\".\") (1:22)\n | filter(Repositories, .Name == \"tabia\")\n | .....................^")
	assert.Nil(reduced)

	reduced, err = github.Reduce(repos, `{ UnExistingFunc(.URL, "stuff") }`)
	assert.Error(err)
	assert.EqualError(err, "cannot get \"UnExistingFunc\" from github.RepositoryFilterEnv (1:24)\n | filter(Repositories, { UnExistingFunc(.URL, \"stuff\") })\n | .......................^")
	assert.Nil(reduced)
}

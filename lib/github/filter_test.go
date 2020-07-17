package github_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
)

func TestReduce(t *testing.T) {
	assert := assert.New(t)

	repos := []github.Repository{
		github.Repository{Name: "tabia", Visibility: github.Public, Owner: "philips-labs"},
		github.Repository{Name: "garo", Visibility: github.Public, Owner: "philips-labs"},
		github.Repository{Name: "dct-notary-admin", Visibility: github.Public, Owner: "philips-labs"},
		github.Repository{Name: "company-draft", Visibility: github.Internal, Owner: "philips-labs"},
		github.Repository{Name: "top-secret", Visibility: github.Private, Owner: "philips-labs"},
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

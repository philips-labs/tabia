package github

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

// RepositoryFilterEnv filter environment for repositories
type RepositoryFilterEnv struct {
	Repositories []Repository
}

// Contains reports wether substring is in s.
func (RepositoryFilterEnv) Contains(s, substring string) bool {
	return strings.Contains(s, substring)
}

// IsPublic indicates if a repository has public visibility.
func (r Repository) IsPublic() bool {
	return r.Visibility == Public
}

// IsInternal indicates if a repository has internal visibility.
func (r Repository) IsInternal() bool {
	return r.Visibility == Internal
}

// IsPrivate indicates if a repository has private visibility.
func (r Repository) IsPrivate() bool {
	return r.Visibility == Private
}

// Reduce filters the repositories based on the given filter
func Reduce(repositories []Repository, filter string) ([]Repository, error) {
	if strings.TrimSpace(filter) == "" {
		return repositories, nil
	}

	program, err := expr.Compile(fmt.Sprintf("filter(Repositories, %s)", filter))
	if err != nil {
		return nil, err
	}

	result, err := expr.Run(program, RepositoryFilterEnv{repositories})
	if err != nil {
		return nil, err
	}
	var repos []Repository
	for _, repo := range result.([]interface{}) {
		repos = append(repos, repo.(Repository))
	}
	return repos, nil
}

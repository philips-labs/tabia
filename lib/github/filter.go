package github

import (
	"fmt"
	"strings"
	"time"

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

// HasTopic indicates if a repository has a given topic.
func (r Repository) HasTopic(topic string) bool {
	for _, t := range r.Topics {
		if strings.ToLower(t.Name) == strings.ToLower(topic) {
			return true
		}
	}

	return false
}

// HasLanguage indicates if a repository has a given language.
func (r Repository) HasLanguage(language string) bool {
	for _, l := range r.Languages {
		if strings.ToLower(l.Name) == strings.ToLower(language) {
			return true
		}
	}

	return false
}

// UpdatedSince indicates if a repository has been updated since the given date.
// Date has to be given in RFC3339 format, e.g. `2006-01-02T15:04:05Z07:00`.
func (r Repository) UpdatedSince(date string) bool {
	return equalOrAfter(r.UpdatedAt, date)
}

// PushedSince indicates if a repository has been pushed since the given date.
// Date has to be given in RFC3339 format, e.g. `2006-01-02T15:04:05Z07:00`.
func (r Repository) PushedSince(date string) bool {
	return equalOrAfter(r.PushedAt, date)
}

// CreatedSince indicates if a repository has been created since the given date.
// Date has to be given in RFC3339 format, e.g. `2006-01-02T15:04:05Z07:00`.
func (r Repository) CreatedSince(date string) bool {
	return equalOrAfter(r.CreatedAt, date)
}

func equalOrAfter(a time.Time, date string) bool {
	since, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return true
	}

	return a.Equal(since) || a.After(since)
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

package grimoirelab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"

	"github.com/philips-labs/tabia/lib/github"
	"github.com/philips-labs/tabia/lib/shared"
)

// GithubMetadataFactory allows to provide a custom generated metadata
type GithubMetadataFactory func(repo github.Repository) Metadata

// GithubProjectMatcher matches a repository with a project
type GithubProjectMatcher struct {
	Rules map[string]GithubProjectMatcherRule `json:"rules,omitempty"`
}

// GithubProjectMatcherRule rule that matches a repository to a project
type GithubProjectMatcherRule struct {
	URL *Regexp `json:"url,omitempty"`
}

// Regexp embeds a regexp.Regexp, and adds Text/JSON
// (un)marshaling.
type Regexp struct {
	regexp.Regexp
}

// Compile wraps the result of the standard library's
// regexp.Compile, for easy (un)marshaling.
func Compile(expr string) (*Regexp, error) {
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &Regexp{*r}, nil
}

// MustCompile wraps the result of the standard library's
// regexp.Compile, for easy (un)marshaling.
func MustCompile(expr string) *Regexp {
	r := regexp.MustCompile(expr)
	return &Regexp{*r}
}

// UnmarshalText satisfies the encoding.TextMarshaler interface,
// also used by json.Unmarshal.
func (r *Regexp) UnmarshalText(b []byte) error {
	rr, err := Compile(string(b))
	if err != nil {
		return err
	}

	*r = *rr

	return nil
}

// MarshalText satisfies the encoding.TextMarshaler interface,
// also used by json.Marshal.
func (r *Regexp) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// NewGithubProjectMatcherFromJSON initializes GithubProjectMatcher from json
func NewGithubProjectMatcherFromJSON(data io.Reader) (*GithubProjectMatcher, error) {
	var m GithubProjectMatcher
	err := json.NewDecoder(data).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, err
}

// ConvertGithubToProjectsJSON converts the repositories into grimoirelab projects.json
func ConvertGithubToProjectsJSON(repos []github.Repository, metadataFactory GithubMetadataFactory, projectMatcher *GithubProjectMatcher) Projects {
	results := make(Projects)
	bbUser := os.Getenv("TABIA_GITHUB_USER")
	bbToken := os.Getenv("TABIA_GITHUB_TOKEN")
	basicAuth := fmt.Sprintf("%s:%s", bbUser, bbToken)
	for _, repo := range repos {
		projectName := getProjectName(repo, projectMatcher)
		project, found := results[projectName]
		if !found {
			results[projectName] = &Project{}
			project = results[projectName]
			project.Git = make([]string, 0)
		}
		updateFromGithubProject(project, repo, basicAuth, metadataFactory)
	}

	return results
}

func getProjectName(repo github.Repository, projectMatcher *GithubProjectMatcher) string {
	if projectMatcher != nil {
		for k, v := range projectMatcher.Rules {
			if v.URL != nil && v.URL.MatchString(repo.URL) {
				return k
			}
		}
	}

	// fallback to github organization name
	return repo.Owner
}

func updateFromGithubProject(project *Project, repo github.Repository, basicAuth string, metadataFactory GithubMetadataFactory) {
	project.Metadata = metadataFactory(repo)
	link := repo.URL
	if link != "" {
		if repo.Visibility != shared.Public {
			u, _ := url.Parse(link)
			link = fmt.Sprintf("%s://%s@%s%s", u.Scheme, basicAuth, u.Hostname(), u.EscapedPath())
		}
		project.Git = append(project.Git, link+".git")
		project.Github = append(project.Github, link)
		project.GithubRepo = append(project.GithubRepo, link)
	}
}

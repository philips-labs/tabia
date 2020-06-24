package grimoirelab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"

	"github.com/philips-labs/tabia/lib/github"
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

// Regexp adds unmarshalling from json for regexp.Regexp
type Regexp struct {
	*regexp.Regexp
}

// UnmarshalText unmarshals json into a regexp.Regexp
func (r *Regexp) UnmarshalText(b []byte) error {
	regex, err := regexp.Compile(string(b))
	if err != nil {
		return err
	}

	r.Regexp = regex

	return nil
}

// MarshalText marshals regexp.Regexp as string
func (r *Regexp) MarshalText() ([]byte, error) {
	if r.Regexp != nil {
		return []byte(r.Regexp.String()), nil
	}

	return nil, nil
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
	return repo.Owner.Login
}

func updateFromGithubProject(project *Project, repo github.Repository, basicAuth string, metadataFactory GithubMetadataFactory) {
	project.Metadata = metadataFactory(repo)
	link := repo.URL
	if link != "" {
		if repo.IsPrivate {
			u, _ := url.Parse(link)
			link = fmt.Sprintf("%s://%s@%s%s", u.Scheme, basicAuth, u.Hostname(), u.EscapedPath())
		}
		project.Git = append(project.Git, link+".git")
		project.Github = append(project.Github, link)
		project.GithubRepo = append(project.GithubRepo, link)
	}
}

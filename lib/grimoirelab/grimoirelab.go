package grimoirelab

import (
	"fmt"
	"net/url"
	"os"

	"github.com/philips-labs/tabia/lib/bitbucket"
	"github.com/philips-labs/tabia/lib/github"
)

// Projects holds all projects to be loaded in Grimoirelab
type Projects map[string]*Project

// Project holds the project resources and metadata
type Project struct {
	Metadata   Metadata `json:"meta,omitempty"`
	Git        []string `json:"git,omitempty"`
	Github     []string `json:"github,omitempty"`
	GithubRepo []string `json:"github:repo,omitempty"`
}

// Metadata hold metadata for a given project
type Metadata map[string]string

// BitbucketMetadataFactory allows to provide a custom generated metadata
type BitbucketMetadataFactory func(repo bitbucket.Repository) Metadata

// GithubMetadataFactory allows to provide a custom generated metadata
type GithubMetadataFactory func(repo github.Repository) Metadata

// ConvertBitbucketToProjectsJSON converts the repositories into grimoirelab projects.json
func ConvertBitbucketToProjectsJSON(repos []bitbucket.Repository, metadataFactory BitbucketMetadataFactory) Projects {
	results := make(Projects)
	bbUser := os.Getenv("TABIA_BITBUCKET_USER")
	bbToken := os.Getenv("TABIA_BITBUCKET_TOKEN")
	basicAuth := fmt.Sprintf("%s:%s", bbUser, bbToken)
	for _, repo := range repos {
		project, found := results[repo.Project.Name]
		if !found {
			results[repo.Project.Name] = &Project{}
			project = results[repo.Project.Name]
			project.Git = make([]string, 0)
		}
		updateFromBitbucketProject(project, repo, basicAuth, metadataFactory)
	}

	return results
}

// ConvertGithubToProjectsJSON converts the repositories into grimoirelab projects.json
func ConvertGithubToProjectsJSON(repos []github.Repository, metadataFactory GithubMetadataFactory) Projects {
	results := make(Projects)
	bbUser := os.Getenv("TABIA_GITHUB_USER")
	bbToken := os.Getenv("TABIA_GITHUB_TOKEN")
	basicAuth := fmt.Sprintf("%s:%s", bbUser, bbToken)
	for _, repo := range repos {
		project, found := results[repo.Owner.Login]
		if !found {
			results[repo.Owner.Login] = &Project{}
			project = results[repo.Owner.Login]
			project.Git = make([]string, 0)
		}
		updateFromGithubProject(project, repo, basicAuth, metadataFactory)
	}

	return results
}

func updateFromBitbucketProject(project *Project, repo bitbucket.Repository, basicAuth string, metadataFactory BitbucketMetadataFactory) {
	project.Metadata = metadataFactory(repo)
	link := getCloneLink(repo, "http")
	if link != "" {
		if !repo.Public {
			u, _ := url.Parse(link)
			link = fmt.Sprintf("%s://%s@%s%s", u.Scheme, basicAuth, u.Hostname(), u.EscapedPath())
		}
		project.Git = append(project.Git, link)
	}
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

func getCloneLink(repo bitbucket.Repository, linkName string) string {
	for _, l := range repo.Links.Clone {
		if l.Name == linkName {
			return l.Href
		}
	}
	return ""
}

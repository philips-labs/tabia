package grimoirelab

import (
	"fmt"
	"net/url"
	"os"

	"github.com/philips-labs/tabia/lib/github"
)

// GithubMetadataFactory allows to provide a custom generated metadata
type GithubMetadataFactory func(repo github.Repository) Metadata

// ConvertGithubToProjectsJSON converts the repositories into grimoirelab projects.json
func ConvertGithubToProjectsJSON(repos []github.Repository, metadataFactory GithubMetadataFactory) Projects {
	results := make(Projects)
	bbUser := os.Getenv("TABIA_GITHUB_USER")
	bbToken := os.Getenv("TABIA_GITHUB_TOKEN")
	basicAuth := fmt.Sprintf("%s:%s", bbUser, bbToken)
	for _, repo := range repos {
		projectName := getProjectName(repo)
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

func getProjectName(repo github.Repository) string {
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

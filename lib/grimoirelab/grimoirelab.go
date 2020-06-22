package grimoirelab

import (
	"fmt"
	"os"
	"strings"

	"github.com/philips-labs/tabia/lib/bitbucket"
)

// Projects holds all projects to be loaded in Grimoirelab
type Projects map[string]Project

// Project holds the project resources and metadata
type Project struct {
	Metadata Metadata `json:"meta,omitempty"`
	Git      []string `json:"git,omitempty"`
}

// Metadata hold metadata for a given project
type Metadata map[string]string

// MetadataFactory allows to provide a custom generated metadata
type MetadataFactory func(repo bitbucket.Repository) Metadata

// ConvertProjectsJSON converts the repositories into grimoirelab projects.json
func ConvertProjectsJSON(repos []bitbucket.Repository, metadataFactory MetadataFactory) Projects {
	results := make(Projects)
	for _, repo := range repos {
		project, found := results[repo.Project.Name]
		if !found {
			results[repo.Project.Name] = Project{}
			project = results[repo.Project.Name]
			project.Git = make([]string, 0)
		}
		updateProject(&project, repo, metadataFactory)
		results[repo.Project.Name] = project
	}

	return results
}

func updateProject(project *Project, repo bitbucket.Repository, metadataFactory MetadataFactory) {
	project.Metadata = metadataFactory(repo)
	link := getCloneLink(repo, "http")
	if link != "" {
		if !repo.Public {
			bbUser := os.Getenv("TABIA_BITBUCKET_USER")
			bbToken := os.Getenv("TABIA_BITBUCKET_TOKEN")
			basicAuth := fmt.Sprintf("%s:%s", bbUser, bbToken)
			url := after(link, "https://")
			link = fmt.Sprintf("https://%s@%s", basicAuth, url)
		}
		project.Git = append(project.Git, link)
	}
}

func after(value, some string) string {
	pos := strings.LastIndex(value, some)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(some)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func getCloneLink(repo bitbucket.Repository, linkName string) string {
	for _, l := range repo.Links.Clone {
		if l.Name == linkName {
			return l.Href
		}
	}
	return ""
}

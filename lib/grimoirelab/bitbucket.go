package grimoirelab

import (
	"fmt"
	"net/url"
	"os"

	"github.com/philips-labs/tabia/lib/bitbucket"
)

// BitbucketMetadataFactory allows to provide a custom generated metadata
type BitbucketMetadataFactory func(repo bitbucket.Repository) Metadata

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

func getCloneLink(repo bitbucket.Repository, linkName string) string {
	for _, l := range repo.Links.Clone {
		if l.Name == linkName {
			return l.Href
		}
	}
	return ""
}

package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/shurcooL/githubv4"

	"github.com/philips-labs/tabia/lib/github/graphql"
	"github.com/philips-labs/tabia/lib/shared"
)

type RestRepo struct {
	Name string
}

type Repository struct {
	ID             string            `json:"id,omitempty"`
	Name           string            `json:"name,omitempty"`
	Description    string            `json:"description,omitempty"`
	URL            string            `json:"url,omitempty"`
	SSHURL         string            `json:"ssh_url,omitempty"`
	Owner          string            `json:"owner,omitempty"`
	Visibility     shared.Visibility `json:"visibility"`
	CreatedAt      time.Time         `json:"created_at,omitempty"`
	UpdatedAt      time.Time         `json:"updated_at,omitempty"`
	PushedAt       time.Time         `json:"pushed_at,omitempty"`
	ForkCount      int               `json:"fork_count,omitempty"`
	StargazerCount int               `json:"stargazer_count,omitempty"`
	WatcherCount   int               `json:"watcher_count,omitempty"`
	Topics         []Topic           `json:"topics,omitempty"`
	Languages      []Language        `json:"languages,omitempty"`
	Collaborators  []Collaborator    `json:"collaborators,omitempty"`
}

type Collaborator struct {
	*graphql.Collaborator
}

type Topic struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Language struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
	Size  int    `json:"size,omitempty"`
}

func (c *Client) FetchOrganziationRepositories(ctx context.Context, owner string) ([]Repository, error) {
	var q struct {
		Repositories graphql.RepositorySearch `graphql:"search(query: $query, type: REPOSITORY, first:100, after: $repoCursor)""`
	}

	// archived repositories are filtered as they give error when fetching collaborators
	// this bug is known with Github.
	// Also see https://github.com/shurcooL/githubv4/issues/72, on proposal for better error handling
	variables := map[string]interface{}{
		"query":      githubv4.String(fmt.Sprintf("org:%s archived:false", owner)),
		"repoCursor": (*githubv4.String)(nil),
	}

	var repositories []graphql.Repository
	for {
		err := c.Query(ctx, &q, variables)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, repositoryEdges(q.Repositories.Edges)...)
		if !q.Repositories.PageInfo.HasNextPage {
			break
		}

		variables["repoCursor"] = githubv4.NewString(q.Repositories.PageInfo.EndCursor)
	}

	// currently the graphql api does not seem to support private vs internal.
	// therefore we use the rest api to fetch the private repos so we can determine private vs internal in the Map function.
	privateRepos, err := c.FetchRestRepositories(ctx, owner, "private")
	if err != nil {
		return nil, err
	}

	return Map(repositories, privateRepos)
}

func (c *Client) FetchRestRepositories(ctx context.Context, owner, repoType string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	opt := &github.RepositoryListByOrgOptions{Type: repoType}
	for {
		repos, resp, err := c.restClient.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

func Map(repositories []graphql.Repository, privateRepositories []*github.Repository) ([]Repository, error) {
	repos := make([]Repository, len(repositories))
	for i, repo := range repositories {
		repos[i] = Repository{
			ID:             repo.ID,
			Name:           repo.Name,
			Description:    strings.TrimSpace(repo.Description),
			URL:            repo.URL,
			SSHURL:         repo.SSHURL,
			Owner:          repo.Owner.Login,
			CreatedAt:      repo.CreatedAt,
			UpdatedAt:      repo.UpdatedAt,
			PushedAt:       repo.PushedAt,
			ForkCount:      repo.ForkCount,
			StargazerCount: repo.Stargazers.TotalCount,
			WatcherCount:   repo.Watchers.TotalCount,
			Topics:         mapTopics(repo.RepositoryTopics),
			Languages:      mapLanguages(repo.Languages),
			Collaborators:  mapCollaborators(repo.Collaborators),
		}

		if repo.IsPrivate {
			isPrivate := false
			for _, privRepo := range privateRepositories {
				if *privRepo.Name == repo.Name {
					isPrivate = true
					break
				}
			}

			if isPrivate {
				repos[i].Visibility = shared.Private
			} else {
				repos[i].Visibility = shared.Internal
			}
		} else {
			repos[i].Visibility = shared.Public
		}
	}

	return repos, nil
}

func mapTopics(topics graphql.RepositoryTopics) []Topic {
	ghTopics := make([]Topic, len(topics.Nodes))
	for i, topic := range topics.Nodes {
		ghTopics[i] = Topic{Name: topic.Topic.Name, URL: fmt.Sprintf("https://github.com%s", topic.ResourcePath)}
	}
	return ghTopics
}

func mapLanguages(languages graphql.Languages) []Language {
	ghLanguages := make([]Language, len(languages.Edges))
	for i, lang := range languages.Edges {
		ghLanguages[i] = Language{Name: lang.Node.Name, Size: lang.Size, Color: lang.Node.Color}
	}
	return ghLanguages
}

func mapCollaborators(collaborators graphql.Collaborators) []Collaborator {
	ghCollaborators := make([]Collaborator, len(collaborators.Nodes))
	for i := range collaborators.Nodes {
		ghCollaborators[i] = Collaborator{&collaborators.Nodes[i]}
	}
	return ghCollaborators
}

func repositoryEdges(edges []graphql.Edge) []graphql.Repository {
	var repositories []graphql.Repository
	for _, edge := range edges {
		repositories = append(repositories, edge.Node.Repository)
	}
	return repositories
}

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/shurcooL/githubv4"

	"github.com/philips-labs/tabia/lib/github/graphql"
)

type Client struct {
	httpClient *http.Client
	*githubv4.Client
}

func NewClientWithTokenAuth(token string) *Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	return &Client{httpClient, client}
}

//go:generate stringer -type=Visibility

// Visibility indicates repository visibility
type Visibility int

const (
	// Public repositories are publicly visible
	Public Visibility = iota
	// Internal repositories are only visible to organization members
	Internal
	// Private repositories are only visible to authorized users
	Private
)

func (v *Visibility) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	default:
	case "public":
		*v = Public
	case "internal":
		*v = Internal
	case "private":
		*v = Private
	}
	return nil
}

func (v Visibility) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

type RestRepo struct {
	Name string
}

type Repository struct {
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	URL         string     `json:"url,omitempty"`
	SSHURL      string     `json:"ssh_url,omitempty"`
	Owner       string     `json:"owner,omitempty"`
	Visibility  Visibility `json:"visibility,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	PushedAt    time.Time  `json:"pushed_at,omitempty"`
	Topics      []Topic    `json:"topics,omitempty"`
}

type Topic struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (c *Client) FetchOrganziationRepositories(ctx context.Context, owner string) ([]Repository, error) {
	var q struct {
		Organization graphql.Organization `graphql:"organization(login: $owner)"`
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repoCursor": (*githubv4.String)(nil),
	}

	var repositories []graphql.Repository
	for {
		err := c.Query(ctx, &q, variables)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, q.Organization.Repositories.Nodes...)
		if !q.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["repoCursor"] = githubv4.NewString(q.Organization.Repositories.PageInfo.EndCursor)
	}

	// currently the graphql api does not seem to support private vs internal.
	// therefore we use the rest api to fetch the private repos so we can determine private vs internal in the Map function.
	privateRepos, err := c.FetchRestRepositories(owner, "private")
	if err != nil {
		return nil, err
	}
	return Map(repositories, privateRepos)
}

func (c *Client) FetchRestRepositories(owner, repoType string) ([]RestRepo, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("https://api.github.com/orgs/%s/repos?type=%s", owner, repoType))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var privateRepos []RestRepo
	err = json.NewDecoder(resp.Body).Decode(&privateRepos)
	if err != nil {
		return nil, err
	}

	return privateRepos, nil
}

func Map(repositories []graphql.Repository, privateRepositories []RestRepo) ([]Repository, error) {
	repos := make([]Repository, len(repositories))
	for i, repo := range repositories {
		repos[i] = Repository{
			ID:          repo.ID,
			Name:        repo.Name,
			Description: strings.TrimSpace(repo.Description),
			URL:         repo.URL,
			SSHURL:      repo.SSHURL,
			Owner:       repo.Owner.Login,
			CreatedAt:   repo.CreatedAt,
			UpdatedAt:   repo.UpdatedAt,
			PushedAt:    repo.PushedAt,
			Topics:      mapTopics(repo.RepositoryTopics),
		}

		if repo.IsPrivate {
			isPrivate := false
			for _, privRepo := range privateRepositories {
				if privRepo.Name == repo.Name {
					isPrivate = true
					break
				}
			}

			if isPrivate {
				repos[i].Visibility = Private
			} else {
				repos[i].Visibility = Internal
			}
		} else {
			repos[i].Visibility = Public
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

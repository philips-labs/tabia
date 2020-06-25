package github

import (
	"context"
	"time"

	"golang.org/x/oauth2"

	"github.com/shurcooL/githubv4"
)

type Client struct {
	*githubv4.Client
}

func NewClientWithTokenAuth(token string) *Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	return &Client{client}
}

type PageInfo struct {
	HasNextPage bool            `json:"has_next_page,omitempty"`
	EndCursor   githubv4.String `json:"end_cursor,omitempty"`
}

type Owner struct {
	Login string `json:"login,omitempty"`
}

type Repository struct {
	Name      string    `json:"name,omitempty"`
	ID        string    `json:"id,omitempty"`
	URL       string    `json:"url,omitempty"`
	SSHURL    string    `json:"ssh_url,omitempty"`
	Owner     Owner     `json:"owner,omitempty"`
	IsPrivate bool      `json:"is_private,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	PushedAt  time.Time `json:"pushed_at,omitempty"`
}

type Repositories struct {
	TotalCount int          `json:"total_count,omitempty"`
	PageInfo   PageInfo     `json:"page_info,omitempty"`
	Nodes      []Repository `json:"nodes,omitempty"`
}

type Organization struct {
	Repositories Repositories `graphql:"repositories(first: 100, after: $repoCursor)" json:"repositories,omitempty"`
}

func (c *Client) FetchOrganziationRepositories(ctx context.Context, owner string) ([]Repository, error) {
	var q struct {
		Organization Organization `graphql:"organization(login: $owner)"`
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repoCursor": (*githubv4.String)(nil),
	}

	var repositories []Repository
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

	return repositories, nil
}

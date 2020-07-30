package graphql

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type PageInfo struct {
	HasNextPage bool            `json:"has_next_page,omitempty"`
	EndCursor   githubv4.String `json:"end_cursor,omitempty"`
}

type Owner struct {
	Login string `json:"login,omitempty"`
}

type Repository struct {
	ID               string           `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	Description      string           `json:"description,omitempty"`
	URL              string           `json:"url,omitempty"`
	SSHURL           string           `json:"ssh_url,omitempty"`
	Owner            Owner            `json:"owner,omitempty"`
	IsPrivate        bool             `json:"is_private,omitempty"`
	CreatedAt        time.Time        `json:"created_at,omitempty"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty"`
	PushedAt         time.Time        `json:"pushed_at,omitempty"`
	RepositoryTopics RepositoryTopics `graphql:"repositoryTopics(first: 25)" json:"repository_topics,omitempty"`
	Collaborators    Collaborators    `graphql:"collaborators(first: 15, affiliation: DIRECT)" json:"collaborators,omitempty"`
}

type RepositoryTopics struct {
	Nodes []RepositoryTopic `json:"nodes,omitempty"`
}

type RepositoryTopic struct {
	Topic        Topic  `json:"topic,omitempty"`
	ResourcePath string `json:"resource_path,omitempty"`
}

type Topic struct {
	Name string `json:"name,omitempty"`
}

type Repositories struct {
	TotalCount int          `json:"total_count,omitempty"`
	PageInfo   PageInfo     `json:"page_info,omitempty"`
	Nodes      []Repository `json:"nodes,omitempty"`
}

type Collaborators struct {
	TotalCount int            `json:"total_count,omitempty"`
	PageInfo   PageInfo       `json:"page_info,omitempty"`
	Nodes      []Collaborator `json:"nodes,omitempty"`
}

type Collaborator struct {
	Name      string `json:"name,omitempty"`
	Login     string `json:"login,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type Organization struct {
	Repositories Repositories `graphql:"repositories(first: 100, after: $repoCursor)" json:"repositories,omitempty"`
}

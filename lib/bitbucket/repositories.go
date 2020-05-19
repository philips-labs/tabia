package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
)

type Repositories struct {
	c *Client
}

type Repository struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ScmID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
	Links         Links   `json:"links"`
}

type RepositoriesResponse struct {
	*PagedResponse
	Values []Repository `json:"values"`
}

func (r *Repositories) List(project string) (*RepositoriesResponse, error) {
	url := fmt.Sprintf("%s/projects/%s/repos?limit=100", r.c.baseEndpoint, project)
	repos, err := r.c.RawRequest("GET", url, "")
	if err != nil {
		return nil, err
	}
	defer repos.Close()

	return decodeRepos(repos)
}

func decodeRepos(body io.ReadCloser) (*RepositoriesResponse, error) {
	var result RepositoriesResponse
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to read repositories json: %w", err)
	}

	return &result, nil
}

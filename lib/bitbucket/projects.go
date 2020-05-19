package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
)

type Projects struct {
	c *Client
}

type Project struct {
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       Links  `json:"links"`
}

type ProjectsResponse struct {
	*PagedResponse
	Values []Project `json:"values"`
}

func (r *Projects) List(start int) (*ProjectsResponse, error) {
	url := fmt.Sprintf("%s/projects?limit=100&start=%d", r.c.baseEndpoint, start)
	body, err := r.c.RawRequest("GET", url, "")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return decodeProjects(body)
}

func decodeProjects(body io.ReadCloser) (*ProjectsResponse, error) {
	var result ProjectsResponse

	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to read projects json: %w", err)
	}

	return &result, nil
}

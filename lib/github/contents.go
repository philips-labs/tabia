package github

import (
	"context"
	"io"
)

// DownloadContents downloads file contents from the given filepath
func (c *Client) DownloadContents(ctx context.Context, owner, repo, filepath string) ([]byte, error) {
	contents, _, err := c.restClient.Repositories.DownloadContents(ctx, owner, repo, filepath, nil)
	if err != nil {
		return nil, err
	}

	defer contents.Close()
	return io.ReadAll(contents)
}

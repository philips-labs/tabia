package github

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v32/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/philips-labs/tabia/lib/transport"
)

type Client struct {
	httpClient *http.Client
	restClient *github.Client
	*githubv4.Client
}

func NewClientWithTokenAuth(token string, writer io.Writer) *Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)
	if writer != nil {
		httpClient.Transport = transport.TeeRoundTripper{
			RoundTripper: httpClient.Transport,
			Writer:       writer,
		}
	}
	client := githubv4.NewClient(httpClient)
	restClient := github.NewClient(httpClient)

	return &Client{httpClient, restClient, client}
}

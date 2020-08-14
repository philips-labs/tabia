package gitlab

import (
	"io"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/xanzy/go-gitlab"

	"github.com/philips-labs/tabia/lib/transport"
)

type Client struct {
	*gitlab.Client
}

func NewClientWithTokenAuth(baseUrl, token string, writer io.Writer) (*Client, error) {
	httpClient := cleanhttp.DefaultPooledClient()
	if writer != nil {
		httpClient.Transport = transport.TeeRoundTripper{
			RoundTripper: httpClient.Transport,
			Writer:       writer,
		}
	}
	c, err := gitlab.NewClient(
		token,
		gitlab.WithHTTPClient(httpClient),
		gitlab.WithBaseURL(baseUrl),
	)
	if err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

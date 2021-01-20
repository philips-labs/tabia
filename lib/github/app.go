package github

import (
	"io"
	"net/http"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"

	"github.com/philips-labs/tabia/lib/transport"
)

// NewClientWithAppAuth creates a new client that authenticates using an app integration ID
// and a app private key
func NewClientWithAppAuth(integrationID int64, privateKey string, writer io.Writer) (*Client, error) {
	config := new(githubapp.Config)
	config.App.IntegrationID = integrationID
	config.App.PrivateKey = privateKey
	config.V3APIURL = "https://api.github.com/"
	config.V4APIURL = "https://api.github.com/graphql"

	cc, err := githubapp.NewDefaultCachingClientCreator(
		*config,
		githubapp.WithClientUserAgent("tabia"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(ClientLogging(writer)),
	)

	client, err := cc.NewAppV4Client()
	if err != nil {
		return nil, err
	}
	restClient, err := cc.NewAppClient()
	if err != nil {
		return nil, err
	}

	return &Client{nil, restClient, client}, nil
}

func ClientLogging(writer io.Writer) githubapp.ClientMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		if writer != nil {
			return transport.TeeRoundTripper{
				RoundTripper: next,
				Writer:       writer,
			}
		}

		return next
	}
}

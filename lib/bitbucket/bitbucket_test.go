package bitbucket_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/philips-labs/tabia/lib/bitbucket"
)

func bitbucketTestClient(handler http.Handler) (*bitbucket.Client, string, func()) {
	s := httptest.NewServer(handler)
	baseUrl := s.Listener.Addr().String()
	apiUrl := "http://" + baseUrl + "/rest/api/1.0"
	token := os.Getenv("TABIA_BITBUCKET_TOKEN")
	bb := bitbucket.NewClientWithTokenAuth(apiUrl, token)
	bb.HttpClient.Transport = &http.Transport{
		DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
			return net.Dial(network, baseUrl)
		},
	}

	return bb, apiUrl, s.Close
}

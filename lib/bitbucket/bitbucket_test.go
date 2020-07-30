package bitbucket_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/bitbucket"
)

func bitbucketTestClient(handler http.Handler) (*bitbucket.Client, string, func()) {
	s := httptest.NewServer(handler)
	baseUrl := s.Listener.Addr().String()
	apiUrl := "http://" + baseUrl + "/rest/api/1.0"
	token := os.Getenv("TABIA_BITBUCKET_TOKEN")
	bb := bitbucket.NewClientWithTokenAuth(apiUrl, token, nil)
	bb.HttpClient.Transport = &http.Transport{
		DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
			return net.Dial(network, baseUrl)
		},
	}

	return bb, apiUrl, s.Close
}

func TestClientWithTokenAuth(t *testing.T) {
	assert := assert.New(t)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	baseUrl := s.Listener.Addr().String()
	apiUrl := "http://" + baseUrl + "/rest/api/1.0"
	token := "asd12bjkhu23uy12iu3hh"
	bb := bitbucket.NewClientWithTokenAuth(apiUrl, token, nil)

	assert.Equal(bitbucket.TokenAuth{Token: token}, bb.Auth)
}

func TestClientWithBasicAuth(t *testing.T) {
	assert := assert.New(t)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	baseUrl := s.Listener.Addr().String()
	apiUrl := "http://" + baseUrl + "/rest/api/1.0"
	user := "johndoe"
	pass := "S3cr3t!"
	bb := bitbucket.NewClientWithBasicAuth(apiUrl, user, pass, nil)

	assert.Equal(bitbucket.BasicAuth{Username: user, Password: pass}, bb.Auth)
}

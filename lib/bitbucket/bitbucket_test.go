package bitbucket_test

import (
	"os"

	"github.com/philips-labs/tabia/lib/bitbucket"
)

var (
	apiBaseURL string
	bb         *bitbucket.Client
)

func init() {
	apiBaseURL = os.Getenv("TABIA_BITBUCKET_API")
	token := os.Getenv("TABIA_BITBUCKET_TOKEN")
	bb = bitbucket.NewClientWithTokenAuth(apiBaseURL, token)
	if proxy := os.Getenv("SOCKS_PROXY"); proxy != "" {
		bb.SetSocksProxy(proxy)
	}
}

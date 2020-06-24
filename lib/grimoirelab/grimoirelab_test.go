package grimoirelab_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertUrlHasBasicAuth(t *testing.T, uri, scheme, user, password, hostname, path string) {
	assert := assert.New(t)
	u, err := url.Parse(uri)
	assert.NoError(err)
	assert.Equal(scheme, u.Scheme)
	assert.Equal(user, u.User.Username())
	pass, isSet := u.User.Password()
	assert.True(isSet)
	assert.Equal(password, pass)
	assert.Equal(hostname, u.Hostname())
	assert.Equal(path, u.EscapedPath())
}

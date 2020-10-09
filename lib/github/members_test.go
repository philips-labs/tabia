package github_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
	"github.com/philips-labs/tabia/lib/github/graphql"
)

func TestMapIdentitiesToMembers(t *testing.T) {
	assert := assert.New(t)

	graphqlIdentities := []graphql.ExternalIdentityNode{
		graphql.ExternalIdentityNode{
			SamlIdentity: graphql.SamlIdentity{NameId: "ldapID1234", Username: "ldapID1234"},
			User: graphql.Member{ID: "githubID1234", Login: "marcofranssen", Name: "Marco Franssen", Organization: struct {
				Name string `json:"name,omitempty"`
			}{Name: "philips-labs"}}},
		graphql.ExternalIdentityNode{
			SamlIdentity: graphql.SamlIdentity{NameId: "ldapID5678", Username: "ldapID5678"},
			User: graphql.Member{ID: "githubID5678", Login: "jdoe", Name: "John Doe", Organization: struct {
				Name string `json:"name,omitempty"`
			}{Name: "philips-labs"}},
		},
	}

	ghMembers := github.MapIdentitiesToMembers(graphqlIdentities)

	assert.Len(ghMembers, 2)
	assert.Equal(graphqlIdentities[0].User.ID, ghMembers[0].ID)
	assert.Equal(graphqlIdentities[1].User.ID, ghMembers[1].ID)

	assert.Equal(graphqlIdentities[0].User.Login, ghMembers[0].Login)
	assert.Equal(graphqlIdentities[1].User.Login, ghMembers[1].Login)

	assert.Equal(graphqlIdentities[0].User.Name, ghMembers[0].Name)
	assert.Equal(graphqlIdentities[1].User.Name, ghMembers[1].Name)

	assert.Equal(graphqlIdentities[0].User.Organization.Name, ghMembers[0].Organization)
	assert.Equal(graphqlIdentities[1].User.Organization.Name, ghMembers[1].Organization)

	assert.Equal(graphqlIdentities[0].SamlIdentity.NameId, ghMembers[0].SamlIdentity.ID)
	assert.Equal(graphqlIdentities[1].SamlIdentity.NameId, ghMembers[1].SamlIdentity.ID)
}

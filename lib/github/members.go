package github

import (
	"context"

	"github.com/shurcooL/githubv4"

	"github.com/philips-labs/tabia/lib/github/graphql"
)

type Member struct {
	ID           string        `json:"id,omitempty"`
	Login        string        `json:"login,omitempty"`
	Name         string        `json:"name,omitempty"`
	Email        string        `json:"email,omitempty"`
	Organization string        `json:"organization,omitempty"`
	SamlIdentity *SamlIdentity `json:"saml_identity,omitempty"`
}

type SamlIdentity struct {
	ID string `json:"id,omitempty"`
}

func (c *Client) FetchOrganziationMembers(ctx context.Context, enterprise, organization string) ([]Member, error) {
	var q struct {
		Enterprise graphql.Enterprise `graphql:"enterprise(slug: $enterprise)"`
	}

	variables := map[string]interface{}{
		"enterprise":     githubv4.String(enterprise),
		"organization":   githubv4.String(organization),
		"identityCursor": (*githubv4.String)(nil),
	}

	var identities []graphql.ExternalIdentityNode
	for {
		err := c.Query(ctx, &q, variables)
		if err != nil {
			return nil, err
		}
		idp := q.Enterprise.OwnerInfo.SamlIdentityProvider
		identities = append(identities, identityEdges(idp.ExternalIdentities.Edges)...)
		if !idp.ExternalIdentities.PageInfo.HasNextPage {
			break
		}

		variables["identityCursor"] = githubv4.NewString(idp.ExternalIdentities.PageInfo.EndCursor)
	}

	return MapIdentitiesToMembers(identities), nil
}

func identityEdges(edges []graphql.ExternalIdentityEdge) []graphql.ExternalIdentityNode {
	var identities []graphql.ExternalIdentityNode
	for _, edge := range edges {
		// seems users without the organization field populated are
		// still returned by the api despite the filter on this field
		if edge.Node.User.Organization.Name != "" {
			identities = append(identities, edge.Node)
		}
	}
	return identities
}

func MapIdentitiesToMembers(identities []graphql.ExternalIdentityNode) []Member {
	members := make([]Member, len(identities))
	for i, identity := range identities {
		members[i] = Member{
			ID:           identity.User.ID,
			Login:        identity.User.Login,
			Name:         identity.User.Name,
			Email:        identity.User.Email,
			Organization: identity.User.Organization.Name,
			SamlIdentity: &SamlIdentity{ID: identity.SamlIdentity.NameId},
		}
	}

	return members
}

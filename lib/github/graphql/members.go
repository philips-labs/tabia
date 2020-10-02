package graphql

type Enterprise struct {
	OwnerInfo OwnerInfo `json:"ownerInfo,omitempty"`
}

type OwnerInfo struct {
	SamlIdentityProvider SamlIdentityProvider `json:"samlIdentityProvider,omitempty"`
}

type SamlIdentityProvider struct {
	SsoURL             string             `json:"ssoUrl,omitempty"`
	ExternalIdentities ExternalIdentities `graphql:"externalIdentities(first: 100, after: $identityCursor)" json:"externalIdentities,omitempty"`
}

type ExternalIdentities struct {
	Edges    []ExternalIdentityEdge `json:"edges,omitempty"`
	PageInfo PageInfo               `json:"pageInfo,omitempty"`
}

type ExternalIdentityEdge struct {
	Node ExternalIdentityNode `json:"node,omitempty"`
}

type ExternalIdentityNode struct {
	SamlIdentity SamlIdentity `json:"samlIdentity,omitempty"`
	User         Member       `json:"user,omitempty"`
}

type SamlIdentity struct {
	NameId   string `json:"nameId,omitempty"`
	Username string `json:"username,omitempty"`
}

type Member struct {
	ID           string           `json:"id,omitempty"`
	Login        string           `json:"login,omitempty"`
	Name         string           `json:"name,omitempty"`
	Organization OrganizationName `graphql:"organization(login: $organization)" json:"organization,omitempty"`
}

type OrganizationName struct {
	Name string `json:"name,omitempty"`
}

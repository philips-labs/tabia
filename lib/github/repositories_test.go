package github_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	gh "github.com/google/go-github/v33/github"
	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/tabia/lib/github"
	"github.com/philips-labs/tabia/lib/github/graphql"
	"github.com/philips-labs/tabia/lib/shared"
)

func TestRepositoryVisibilityToJSON(t *testing.T) {
	assert := assert.New(t)

	expectedTemplate := `{"name":"%s","visibility":"%s","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","pushed_at":"0001-01-01T00:00:00Z"}
`

	var result strings.Builder
	jsonEnc := json.NewEncoder(&result)

	privRepo := github.Repository{
		Name:       "private-repo",
		Visibility: shared.Private,
	}
	err := jsonEnc.Encode(privRepo)
	assert.NoError(err)
	assert.Equal(fmt.Sprintf(expectedTemplate, "private-repo", "Private"), result.String())

	var unmarshalledRepo github.Repository
	err = json.Unmarshal([]byte(result.String()), &unmarshalledRepo)
	if assert.NoError(err) {
		assert.Equal(shared.Private, unmarshalledRepo.Visibility)
	}

	internalRepo := github.Repository{
		Name:       "internal-repo",
		Visibility: shared.Internal,
	}
	result.Reset()
	err = jsonEnc.Encode(internalRepo)
	assert.NoError(err)
	assert.Equal(fmt.Sprintf(expectedTemplate, "internal-repo", "Internal"), result.String())

	err = json.Unmarshal([]byte(result.String()), &unmarshalledRepo)
	if assert.NoError(err) {
		assert.Equal(shared.Internal, unmarshalledRepo.Visibility)
	}

	publicRepo := github.Repository{
		Name:       "public-repo",
		Visibility: shared.Public,
	}
	result.Reset()

	err = jsonEnc.Encode(publicRepo)
	assert.NoError(err)
	assert.Equal(fmt.Sprintf(expectedTemplate, "public-repo", "Public"), result.String())

	err = json.Unmarshal([]byte(result.String()), &unmarshalledRepo)
	if assert.NoError(err) {
		assert.Equal(shared.Public, unmarshalledRepo.Visibility)
	}
}

func TestMap(t *testing.T) {
	assert := assert.New(t)

	owner := graphql.Owner{Login: "philips-labs"}
	topics := graphql.RepositoryTopics{
		Nodes: []graphql.RepositoryTopic{
			graphql.RepositoryTopic{Topic: graphql.Topic{Name: "opensource"}, ResourcePath: "/topics/opensource"},
			graphql.RepositoryTopic{Topic: graphql.Topic{Name: "golang"}, ResourcePath: "/topics/golang"},
			graphql.RepositoryTopic{Topic: graphql.Topic{Name: "graphql"}, ResourcePath: "/topics/graphql"},
		},
	}
	collaborators := graphql.Collaborators{
		Nodes: []graphql.Collaborator{
			graphql.Collaborator{Name: "Marco Franssen", Login: "marcofranssen", AvatarURL: "https://avatars3.githubusercontent.com/u/694733?u=6aeb327c48cb88ae31eb88e680b96228f53cae51&v=4"},
			graphql.Collaborator{Name: "John Doe", Login: "johndoe", AvatarURL: "https://avatars3.githubusercontent.com/u/694733?u=6aeb327c48cb88ae31eb88e680b96228f53cae51&v=4"},
		},
	}
	languages := graphql.Languages{
		Edges: []graphql.LanguageEdge{
			graphql.LanguageEdge{Node: graphql.LanguageNode{Name: "Go", Color: "#cc0000"}, Size: 3000},
			graphql.LanguageEdge{Node: graphql.LanguageNode{Name: "JavaScript", Color: "#0000cc"}, Size: 532},
		},
	}
	graphqlRepositories := []graphql.Repository{
		graphql.Repository{Owner: owner, Name: "private-repo", Description: "I am private ", IsPrivate: true},
		graphql.Repository{Owner: owner, Name: "internal-repo", Description: "Superb inner-source stuff", IsPrivate: true},
		graphql.Repository{Owner: owner, Name: "opensource", Description: "I'm shared with the world", RepositoryTopics: topics},
		graphql.Repository{Owner: owner, Name: "secret-repo", Description: " ** secrets ** ", IsPrivate: true, Collaborators: collaborators, Languages: languages},
	}

	privateRepos := []*gh.Repository{
		&gh.Repository{Name: stringPointer("private-repo")},
		&gh.Repository{Name: stringPointer("secret-repo")},
	}
	ghRepos, err := github.Map(graphqlRepositories, privateRepos)
	if !assert.NoError(err) {
		return
	}

	assert.Len(ghRepos, 4)
	assert.Equal(shared.Private, ghRepos[0].Visibility)
	assert.Equal(shared.Internal, ghRepos[1].Visibility)
	assert.Equal(shared.Public, ghRepos[2].Visibility)
	assert.Equal(shared.Private, ghRepos[3].Visibility)

	assert.Equal(owner.Login, ghRepos[0].Owner)
	assert.Equal(owner.Login, ghRepos[1].Owner)
	assert.Equal(owner.Login, ghRepos[2].Owner)
	assert.Equal(owner.Login, ghRepos[3].Owner)

	assert.Equal("I am private", ghRepos[0].Description)
	assert.Equal("Superb inner-source stuff", ghRepos[1].Description)
	assert.Equal("I'm shared with the world", ghRepos[2].Description)
	assert.Equal("** secrets **", ghRepos[3].Description)

	assert.Equal("opensource", ghRepos[2].Topics[0].Name)
	assert.Equal("https://github.com/topics/opensource", ghRepos[2].Topics[0].URL)
	assert.Equal("golang", ghRepos[2].Topics[1].Name)
	assert.Equal("https://github.com/topics/golang", ghRepos[2].Topics[1].URL)
	assert.Equal("graphql", ghRepos[2].Topics[2].Name)
	assert.Equal("https://github.com/topics/graphql", ghRepos[2].Topics[2].URL)

	assert.Len(ghRepos[3].Collaborators, 2)
	assert.Equal("Marco Franssen", ghRepos[3].Collaborators[0].Name)
	assert.Equal("marcofranssen", ghRepos[3].Collaborators[0].Login)
	assert.Equal("https://avatars3.githubusercontent.com/u/694733?u=6aeb327c48cb88ae31eb88e680b96228f53cae51&v=4", ghRepos[3].Collaborators[0].AvatarURL)

	assert.Equal("John Doe", ghRepos[3].Collaborators[1].Name)
	assert.Equal("johndoe", ghRepos[3].Collaborators[1].Login)
	assert.Equal("https://avatars3.githubusercontent.com/u/694733?u=6aeb327c48cb88ae31eb88e680b96228f53cae51&v=4", ghRepos[3].Collaborators[1].AvatarURL)

	assert.Len(ghRepos[3].Languages, 2)
	assert.Equal("Go", ghRepos[3].Languages[0].Name)
	assert.Equal(3000, ghRepos[3].Languages[0].Size)
	assert.Equal("#cc0000", ghRepos[3].Languages[0].Color)
	assert.Equal("JavaScript", ghRepos[3].Languages[1].Name)
	assert.Equal(532, ghRepos[3].Languages[1].Size)
	assert.Equal("#0000cc", ghRepos[3].Languages[1].Color)
}

func stringPointer(s string) *string {
	return &s
}

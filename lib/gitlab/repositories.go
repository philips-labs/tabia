package gitlab

import (
	"context"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"

	"github.com/philips-labs/tabia/lib/shared"
)

type Repository struct {
	ID             int               `json:"id,omitempty"`
	Name           string            `json:"name,omitempty"`
	Description    string            `json:"description,omitempty"`
	Owner          string            `json:"owner,omitempty"`
	URL            string            `json:"url,omitempty"`
	SSHURL         string            `json:"sshurl,omitempty"`
	CreatedAt      *time.Time        `json:"created_at,omitempty"`
	LastActivityAt *time.Time        `json:"last_activity_at,omitempty"`
	Visibility     shared.Visibility `json:"visibility,omitempty"`
}

func (c *Client) ListRepositories(ctx context.Context) ([]Repository, error) {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	var repos []Repository
	for {
		projects, resp, err := c.Projects.ListProjects(opt, gitlab.WithContext(ctx))
		if err != nil {
			return repos, err
		}
		defer resp.Body.Close()
		repos = append(repos, Map(projects)...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}
	return repos, nil
}

func Map(projects []*gitlab.Project) []Repository {
	repos := make([]Repository, len(projects))
	for i, project := range projects {
		repos[i] = Repository{
			ID:             project.ID,
			Name:           project.Name,
			Description:    strings.TrimSpace(project.Description),
			Owner:          mapOwner(project.Owner),
			URL:            project.WebURL,
			SSHURL:         project.SSHURLToRepo,
			CreatedAt:      project.CreatedAt,
			LastActivityAt: project.LastActivityAt,
			Visibility:     shared.VisibilityFromText(string(project.Visibility)),
		}
	}
	return repos
}

func mapOwner(owner *gitlab.User) string {
	if owner != nil {
		return owner.Name
	}
	return ""
}

func boolPointer(b bool) *bool {
	return &b
}

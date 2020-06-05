package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/philips-labs/tabia/lib/bitbucket"
	"github.com/urfave/cli/v2"
)

func createBitbucket() *cli.Command {
	return &cli.Command{
		Name:  "bitbucket",
		Usage: "Gets you some insight in Bitbucket repositories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "api",
				Usage:       "The api enpoint `ENDPOINT`",
				DefaultText: "https://bitbucket.atlas.philips.com/rest/api/1.0",
				EnvVars:     []string{"TABIA_BITBUCKET_API"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "Calls the api using the given `TOKEN`",
				EnvVars:  []string{"TABIA_BITBUCKET_TOKEN"},
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:   "projects",
				Usage:  "display insights on projects",
				Action: projects,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "json",
						Usage: "outputs results in JSON format",
					},
				},
			},
			{
				Name:   "repositories",
				Usage:  "display insights on repositories",
				Action: repositories,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "projects",
						Aliases:  []string{"P"},
						Usage:    "fetches repositories for given projects",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "json",
						Usage: "outputs results in JSON format",
					},
				},
			},
		},
	}
}

func projects(c *cli.Context) error {
	api := c.String("api")
	token := c.String("token")
	asJSON := c.Bool("json")

	bb := bitbucket.NewClientWithTokenAuth(api, token)
	projects := make([]bitbucket.Project, 0)
	page := 0
	for {
		resp, err := bb.Projects.List(page)
		if err != nil {
			return err
		}
		projects = append(projects, resp.Values...)
		page = resp.NextPageStart
		if resp.IsLastPage {
			break
		}
	}

	if asJSON {
		err := printJSON(c.App.Writer, projects)
		if err != nil {
			return err
		}
	} else {
		w := tabwriter.NewWriter(c.App.Writer, 3, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "ID", "Key", "Name", "Public")
		for _, project := range projects {
			fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", project.ID, project.Key, project.Name, project.Public)
		}
		w.Flush()
	}

	return nil
}

func repositories(c *cli.Context) error {
	api := c.String("api")
	token := c.String("token")
	asJSON := c.Bool("json")
	projects := c.StringSlice("projects")

	bb := bitbucket.NewClientWithTokenAuth(api, token)

	results := make([]bitbucket.Repository, 0)
	for _, project := range projects {
		resp, err := bb.Repositories.List(project)
		if err != nil {
			return err
		}
		results = append(results, resp.Values...)
	}

	if asJSON {
		err := printJSON(c.App.Writer, results)
		if err != nil {
			return err
		}
	} else {
		w := tabwriter.NewWriter(c.App.Writer, 3, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "Project\tID\tSlug\tName\tPublic\tClone")
		for _, repo := range results {
			httpClone := getCloneURL(repo.Links.Clone, "http")
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%t\t%s\n", repo.Project.Key, repo.ID, repo.Slug, repo.Name, repo.Public, httpClone)
		}
		w.Flush()
	}

	return nil
}

func getCloneURL(links []bitbucket.CloneLink, linkName string) string {
	for _, l := range links {
		if l.Name == linkName {
			return l.Href
		}
	}
	return ""
}

func printJSON(w io.Writer, data interface{}) error {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", json)
	return nil
}

package cmd

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"

	"github.com/philips-labs/tabia/lib/bitbucket"
	"github.com/philips-labs/tabia/lib/grimoirelab"
	"github.com/philips-labs/tabia/lib/output"
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
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Adds verbose logging",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:   "projects",
				Usage:  "display insights on projects",
				Action: bitbucketProjects,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"F"},
						Usage:       "Formats output in the given `FORMAT`",
						EnvVars:     []string{"TABIA_OUTPUT_FORMAT"},
						DefaultText: "",
					}},
			},
			{
				Name:   "repositories",
				Usage:  "display insights on repositories",
				Action: bitbucketRepositories,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Usage: "fetches repositories for all projects",
					},
					&cli.StringSliceFlag{
						Name:    "projects",
						Aliases: []string{"P"},
						Usage:   "fetches repositories for given projects",
					},
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"F"},
						Usage:       "Formats output in the given `FORMAT`",
						EnvVars:     []string{"TABIA_OUTPUT_FORMAT"},
						DefaultText: "",
					},
					&cli.PathFlag{
						Name:      "template",
						Aliases:   []string{"T"},
						Usage:     "Formats output using the given `TEMPLATE`",
						TakesFile: true,
					},
				},
			},
		},
	}
}

func newBitbucketClient(c *cli.Context) *bitbucket.Client {
	api := c.String("api")
	verbose := c.Bool("verbose")
	token := c.String("token")

	var bbWriter io.Writer
	if verbose {
		bbWriter = c.App.Writer
	}
	return bitbucket.NewClientWithTokenAuth(api, token, bbWriter)
}

func bitbucketProjects(c *cli.Context) error {
	format := c.String("format")

	bb := newBitbucketClient(c)
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

	switch format {
	case "json":
		err := output.PrintJSON(c.App.Writer, projects)
		if err != nil {
			return err
		}
	default:
		w := tabwriter.NewWriter(c.App.Writer, 3, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "ID", "Key", "Name", "Public")
		for _, project := range projects {
			fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", project.ID, project.Key, project.Name, project.Public)
		}
		w.Flush()
	}

	return nil
}

func bitbucketRepositories(c *cli.Context) error {
	format := c.String("format")
	projects := c.StringSlice("projects")

	bb := newBitbucketClient(c)

	results := make([]bitbucket.Repository, 0)
	for _, project := range projects {
		resp, err := bb.Repositories.List(project)
		if err != nil {
			return err
		}
		results = append(results, resp.Values...)
	}

	switch format {
	case "json":
		err := output.PrintJSON(c.App.Writer, results)
		if err != nil {
			return err
		}
	case "grimoirelab":
		projects := grimoirelab.ConvertBitbucketToProjectsJSON(results, func(repo bitbucket.Repository) grimoirelab.Metadata {
			return grimoirelab.Metadata{
				"title":   repo.Project.Name,
				"program": "One Codebase",
			}
		})
		err := output.PrintJSON(c.App.Writer, projects)
		if err != nil {
			return err
		}
	case "templated":
		if !c.IsSet("template") {
			return fmt.Errorf("you must specify the path to the template")
		}

		templateFile := c.Path("template")
		tmplContent, err := os.ReadFile(templateFile)
		if err != nil {
			return err
		}
		tmpl, err := template.New("repositories").Parse(string(tmplContent))
		if err != nil {
			return err
		}
		err = tmpl.Execute(c.App.Writer, results)
		if err != nil {
			return err
		}
	default:
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

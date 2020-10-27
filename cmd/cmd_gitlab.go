package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"text/tabwriter"

	"github.com/urfave/cli/v2"

	"github.com/philips-labs/tabia/lib/gitlab"
	"github.com/philips-labs/tabia/lib/output"
)

func createGitlab() *cli.Command {
	return &cli.Command{
		Name:  "gitlab",
		Usage: "Gets you some insight in Gitlab repositories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "Calls the api using the given `TOKEN`",
				EnvVars:  []string{"TABIA_GITLAB_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:        "instance",
				Usage:       "The instance url `URL`",
				DefaultText: "https://gitlab.com/",
				EnvVars:     []string{"TABIA_GITLAB_INSTANCE"},
				Required:    true,
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Adds verbose logging",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:   "repositories",
				Usage:  "display insights on repositories",
				Action: gitlabRepositories,
				Flags: []cli.Flag{
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
					&cli.StringFlag{
						Name:    "filter",
						Aliases: []string{"f"},
						Usage:   "filters repositories based on the given `EXPRESSION`",
					},
				},
			},
		},
	}
}

func newGitlabClient(c *cli.Context) (*gitlab.Client, error) {
	instance := c.String("instance")
	verbose := c.Bool("verbose")
	token := c.String("token")

	var ghWriter io.Writer
	if verbose {
		ghWriter = c.App.Writer
	}

	return gitlab.NewClientWithTokenAuth(instance, token, ghWriter)
}

func gitlabRepositories(c *cli.Context) error {
	format := c.String("format")
	filter := c.String("filter")

	client, err := newGitlabClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	filters := gitlab.ConvertFiltersToListProjectOptions(filter)
	repos, err := client.ListRepositories(ctx, filters...)
	if err != nil {
		return err
	}

	filtered, err := gitlab.Reduce(repos, filter)
	if err != nil {
		return err
	}

	switch format {
	case "json":
		output.PrintJSON(c.App.Writer, filtered)
	case "templated":
		if !c.IsSet("template") {
			return fmt.Errorf("you must specify the path to the template")
		}

		templateFile := c.Path("template")
		tmplContent, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return err
		}
		err = output.PrintUsingTemplate(c.App.Writer, tmplContent, filtered)
		if err != nil {
			return err
		}
	default:
		w := tabwriter.NewWriter(c.App.Writer, 3, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, " \tID\tOwner\tName\tVisibility\tURL")
		for i, repo := range filtered {
			fmt.Fprintf(w, "%04d\t%d\t%s\t%s\t%s\t%s\n", i+1, repo.ID, repo.Owner, repo.Name, repo.Visibility, repo.URL)
		}
		w.Flush()
	}

	return nil
}

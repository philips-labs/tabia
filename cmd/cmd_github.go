package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/philips-labs/tabia/lib/output"

	"github.com/philips-labs/tabia/lib/github"

	"github.com/urfave/cli/v2"
)

func createGithub() *cli.Command {
	return &cli.Command{
		Name:  "github",
		Usage: "Gets you some insight in Github repositories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "Calls the api using the given `TOKEN`",
				EnvVars:  []string{"TABIA_GITHUB_TOKEN"},
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:   "repositories",
				Usage:  "display insights on repositories",
				Action: githubRepositories,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "owner",
						Aliases: []string{"O"},
						Usage:   "fetches repositories for given owner",
					},
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"F"},
						Usage:       "Formats output in the given `FORMAT`",
						EnvVars:     []string{"TABIA_OUTPUT_FORMAT"},
						DefaultText: "",
					},
				},
			},
		},
	}
}

func githubRepositories(c *cli.Context) error {
	owners := c.StringSlice("owner")
	format := c.String("format")

	client := github.NewClientWithTokenAuth(os.Getenv("TABIA_GITHUB_TOKEN"))
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	var repositories []github.Repository
	for _, owner := range owners {
		repos, err := client.FetchOrganziationRepositories(ctx, owner)
		if err != nil {
			return err
		}
		repositories = append(repositories, repos...)
	}

	switch format {
	case "json":
		output.PrintJSON(c.App.Writer, repositories)
	default:
		w := tabwriter.NewWriter(c.App.Writer, 3, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, " \tName\tOwner\tPublic\tClone")
		for i, repo := range repositories {
			fmt.Fprintf(w, "%04d\t%s\t%s\t%t\t%s\n", i+1, repo.Name, repo.Owner.Login, !repo.IsPrivate, repo.URL)
		}
		w.Flush()
	}

	return nil
}

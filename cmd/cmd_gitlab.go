package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
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
			},
		},
	}
}

func gitlabRepositories(c *cli.Context) error {
	instance := c.String("instance")
	token := c.String("token")

	fmt.Fprintln(c.App.Writer, "Not implemented yet.")
	fmt.Fprintf(c.App.Writer, "token: %s\n", token)
	fmt.Fprintf(c.App.Writer, "instance: %s\n", instance)

	return nil
}

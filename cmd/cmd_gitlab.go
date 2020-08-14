package cmd

import (
	"fmt"
	"io"

	"github.com/urfave/cli/v2"

	"github.com/philips-labs/tabia/lib/gitlab"
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
	client, err := newGitlabClient(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.App.Writer, client.BaseURL())

	return nil
}

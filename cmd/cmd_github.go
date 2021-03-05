package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v2"

	"github.com/philips-labs/tabia/lib/github"
	"github.com/philips-labs/tabia/lib/grimoirelab"
	"github.com/philips-labs/tabia/lib/output"
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
				Required: false,
			},
			&cli.StringFlag{
				Name:     "integration-id",
				Aliases:  []string{"app-id"},
				Usage:    "Authenticates to Github using the given `APP_INTEGRATION_ID`",
				EnvVars:  []string{"TABIA_GITHUB_APP_INTEGRATION_ID"},
				Required: false,
			},
			&cli.PathFlag{
				Name:      "private-key",
				Usage:     "Authenticates to Github using the given `APP_PRIVATE_KEY`",
				EnvVars:   []string{"TABIA_GITHUB_APP_PRIVATE_KEY"},
				Required:  false,
				TakesFile: true,
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
				Action: githubRepositories,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "owner",
						Aliases: []string{"O"},
						Usage:   "fetches repositories for given owner",
					},
					&cli.PathFlag{
						Name:      "matching",
						Aliases:   []string{"M"},
						Usage:     "matches repositories to projects based on json file",
						TakesFile: true,
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
					&cli.StringFlag{
						Name:    "filter",
						Aliases: []string{"f"},
						Usage:   "filters repositories based on the given `EXPRESSION`",
					},
				},
			},
			{
				Name:   "members",
				Usage:  "display insights on members",
				Action: githubMembers,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "organization",
						Aliases: []string{"O"},
						Usage:   "fetches members for given organization",
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
					&cli.StringFlag{
						Name:    "filter",
						Aliases: []string{"f"},
						Usage:   "filters members based on the given `EXPRESSION`",
					},
				},
			},
			{
				Name:      "contents",
				Usage:     "Gets contents from a repository",
				Action:    githubContents,
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "repo",
						Aliases:  []string{"R"},
						Usage:    "fetches content of given `REPO`",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "fetches content of given `FILE`",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "writes contents to `FILEPATH`",
						Required: false,
					},
				},
			},
		},
	}
}

func newGithubClient(c *cli.Context) (*github.Client, error) {
	verbose := c.Bool("verbose")

	var ghWriter io.Writer
	if verbose {
		ghWriter = c.App.Writer
	}

	if c.IsSet("integration-id") && c.IsSet("private-key") {
		integrationID := c.Int64("integration-id")
		privateKey := c.Path("private-key")

		privateKeyBytes, err := os.ReadFile(privateKey)
		if err != nil {
			return nil, err
		}
		org := append(c.StringSlice("owner"), c.StringSlice("organization")...)

		client, err := github.NewClientWithAppAuth(integrationID, string(privateKeyBytes), org[0], ghWriter)
		return client, nil
	}

	if !c.IsSet("token") {
		return nil, fmt.Errorf("no `integration-id` and `private-key` or `token` provided")
	}

	token := c.String("token")
	return github.NewClientWithTokenAuth(token, ghWriter), nil
}

func githubMembers(c *cli.Context) error {
	owners := c.StringSlice("organization")
	format := c.String("format")
	filter := c.String("filter")

	client, err := newGithubClient(c)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	var ghMembers []github.Member
	for _, owner := range owners {
		members, err := client.FetchOrganziationMembers(ctx, "royal-philips", owner)
		if err != nil {
			return err
		}
		filtered, err := github.ReduceMembers(members, filter)
		if err != nil {
			return err
		}
		ghMembers = append(ghMembers, filtered...)
	}

	switch format {
	case "json":
		output.PrintJSON(c.App.Writer, ghMembers)
	case "templated":
		if !c.IsSet("template") {
			return fmt.Errorf("you must specify the path to the template")
		}

		templateFile := c.Path("template")
		tmplContent, err := os.ReadFile(templateFile)
		if err != nil {
			return err
		}
		err = output.PrintUsingTemplate(c.App.Writer, tmplContent, ghMembers)
		if err != nil {
			return err
		}
	default:
		w := tabwriter.NewWriter(c.App.Writer, 3, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, " \tLogin\tSAML Identity\tOrganization\tName")
		for i, m := range ghMembers {
			fmt.Fprintf(w, "%04d\t%s\t%s\t%s\t%s\n", i+1, m.Login, m.SamlIdentity.ID, m.Organization, m.Name)
		}
		w.Flush()
	}

	return nil
}

func githubRepositories(c *cli.Context) error {
	owners := c.StringSlice("owner")
	format := c.String("format")
	filter := c.String("filter")

	client, err := newGithubClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	var repositories []github.Repository
	for _, owner := range owners {
		repos, err := client.FetchOrganziationRepositories(ctx, owner)
		if err != nil {
			return err
		}
		filtered, err := github.ReduceRepositories(repos, filter)
		if err != nil {
			return err
		}
		repositories = append(repositories, filtered...)
	}

	switch format {
	case "json":
		output.PrintJSON(c.App.Writer, repositories)
	case "grimoirelab":
		projectMatchingConfig := c.Path("matching")
		json, err := os.Open(projectMatchingConfig)
		if err != nil {
			return err
		}
		defer json.Close()
		projectMatcher, err := grimoirelab.NewGithubProjectMatcherFromJSON(json)
		if err != nil {
			return err
		}

		projects := grimoirelab.ConvertGithubToProjectsJSON(
			repositories,
			func(repo github.Repository) grimoirelab.Metadata {
				return grimoirelab.Metadata{
					"title":   repo.Owner,
					"program": "One Codebase",
				}
			},
			projectMatcher)
		err = output.PrintJSON(c.App.Writer, projects)
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
		err = output.PrintUsingTemplate(c.App.Writer, tmplContent, repositories)
		if err != nil {
			return err
		}
	default:
		w := tabwriter.NewWriter(c.App.Writer, 3, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, " \tName\tOwner\tVisibility\tClone")
		for i, repo := range repositories {
			fmt.Fprintf(w, "%04d\t%s\t%s\t%s\t%s\n", i+1, repo.Name, repo.Owner, repo.Visibility, repo.URL)
		}
		w.Flush()
	}

	return nil
}

func githubContents(c *cli.Context) error {
	repo := c.String("repo")
	filePath := c.String("file")
	output := c.Path("output")

	client, err := newGithubClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	repoParts := strings.Split(repo, "/")
	owner := repoParts[0]
	repo = repoParts[1]

	contents, err := client.DownloadContents(ctx, owner, repo, filePath)
	if err != nil {
		return err
	}

	if output != "" {
		if strings.Contains(output, "/") {
			dir := filepath.Dir(output)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}
		err := os.WriteFile(output, contents, 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Fprintf(c.App.Writer, "%s", contents)
	}

	return nil
}

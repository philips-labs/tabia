package cmd

import "github.com/urfave/cli/v2"

// CreateCommands Creates CLI commands.
func CreateCommands() []*cli.Command {
	return []*cli.Command{
		createBitbucket(),
	}
}

package main

import (
	"os"

	"github.com/lewispeckover/consulator/command"

	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Serf commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"dump": func() (cli.Command, error) {
			return &command.DumpCommand{
				Ui: ui,
			}, nil
		},
		"import": func() (cli.Command, error) {
			return &command.ImportCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Version:   Version,
				BuildDate: BuildDate,
				Ui:        ui,
			}, nil
		},
	}
}

// Initializes the Commands struct for the application.
// New commands need to be added to the CommandFactory map.
package main

import (
	"os"

	"github.com/mohae/cli"
	"github.com/mohae/pipeline/app"
	"github.com/mohae/pipeline/command"
)

// Commands
var Commands map[string]cli.CommandFactory

// Set-up the commands for the application. Help and version doesn't need to bo
// set-up because they are always available.
func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	Commands = map[string]cli.CommandFactory{
		"square": func() (cli.Command, error) {
			return &command.SquareCommand{
				UI: ui,
			}, nil
		},
		"md5": func() (cli.Command, error) {
			return &command.MD5Command{
				UI: ui,
			}, nil
		},
		"sha256": func() (cli.Command, error) {
			return &command.SHA256Command{
				UI: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Name:              app.Name,
				Revision:          GitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
				UI:                ui,
			}, nil
		},
	}
}

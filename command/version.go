package command

import (
	"bytes"
	"fmt"

	"github.com/mohae/cli"
)

// VersionCommand is a Command implementation that prints the version.
type VersionCommand struct {
	Name		string
	Revision          string
	Version           string
	VersionPrerelease string
	UI                cli.Ui
}

// Help prints the Help text for the version sub-command
func (c *VersionCommand) Help() string {
	return "Prints " + c.Name + "'s version information."
}

// Run runs the version sub-command.
func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "%s v%s", c.Name, c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, ".%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	c.UI.Output(versionString.String())

	return 0
}

// Synopsis provides a precis of the version sub-command.
func (c *VersionCommand) Synopsis() string {
	return "Prints the " + c.Name + " version"
}

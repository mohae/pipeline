package command

import (
	"strings"

	"github.com/mohae/cli"
	"github.com/mohae/contour"
	"github.com/mohae/pipeline/app"
)

// SqueareCommand is a Command implementation that says hello world
type SquareCommand struct {
	UI cli.Ui
}

// Help prints the help text for the run sub-command.
func (c *SquareCommand) Help() string {
	helpText := `
Usage: pipeline squaer [flags] <nums int...>

square will take a 1 or more ints and provide their square

    $ quine square 2
	4

    $ quine square 2 3 4
	4
	9
	16
 
`
	return strings.TrimSpace(helpText)
}

// Run runs the square command; the args are a variadic list of strings that
// represent ints to square
func (c *SquareCommand) Run(args []string) int {
	// set up the command flags
	contour.SetFlagSetUsage(func() { c.UI.Output(c.Help()) })

	// Filter the flags from the args and update the config with them.
	// The args remaining after being filtered are returned.
	filteredArgs, err := contour.FilterArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Run the command in the package.
	message, err := app.Square(filteredArgs...)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.UI.Output(message)
	return 0
}

// Synopsis provides a precis of the hello command.
func (c *SquareCommand) Synopsis() string {
	ret := `Computes the square of the provided ints
`
	return ret
}

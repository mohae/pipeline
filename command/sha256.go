package command

import (
	"strings"

	"github.com/mohae/cli"
	"github.com/mohae/contour"
	"github.com/mohae/pipeline/app"
)

// SHA256Command is a Command implementation that says hello world
type SHA256Command struct {
	UI cli.Ui
}

// Help prints the help text for the run sub-command.
func (c *SHA256Command) Help() string {
	helpText := `
Usage: pipeline sha256 [flags] <path/to/dir>

sha256 will take a path, representing either a file or a directory,
calculate the sha256 checksums, and print them out

    $ pipeline sha256 somefile.txt
	d47c2bbc28298ca9befdfbc5d3aa4e65 somefile.txt

    $ pipeline sha256 dirpath
	d47c2bbc28298ca9befdfbc5d3aa4e65 dirpath/somefile.txt
	.
	.
 	.
	
supported flags:

    --parallel=n       set the number of parallel digesters to n
    -p                 alias to --parallel
`
	return strings.TrimSpace(helpText)
}

// Run runs the square command; the args are a variadic list of strings that
// represent ints to square
func (c *SHA256Command) Run(args []string) int {
	// set up the command flags
	contour.SetFlagSetUsage(func() { c.UI.Output(c.Help()) })

	// Filter the flags from the args and update the config with them.
	// The args remaining after being filtered are returned.
	filteredArgs, err := contour.FilterArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Since we only care about the first element, that is all that gets passed.
	// All else gets filtered.
	message, err := app.SHA256(filteredArgs[0])
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.UI.Output(message)
	return 0
}

// Synopsis provides a precis of the hello command.
func (c *SHA256Command) Synopsis() string {
	ret := `Computes the SHA256 hash for files
`
	return ret
}

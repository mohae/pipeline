// Copyright © 2014, All rights reserved
// Joel Scoble, https://github.com/mohae/clitpl
//
// This is licensed under The MIT License. Please refer to the included
// LICENSE file for more information. If the LICENSE file has not been
// included, please refer to the url above.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License
//
// clitpl is a basic implementation of michellh's cli package. clitpl uses
// cli's example, and the implementations used in both mitchellh's and
// hashicorp's implemenations.
//
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mohae/cli"
	"github.com/mohae/pipeline/app"
)

// This is modeled on mitchellh's realmain wrapper
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Exit(realMain())
}

// realMain, is the actual main for the application. This keeps all changes
// needed for a new application to one file in the main application directory.
// In addition to this, only commands/ needs to be modified, adding the app's
// commands and any handler codes for those commands, like the 'cmd' package.
//
// No logging is done until the flags are processed, since the flags could
// enable/disable output, alter it, or alter its output locations. Everything
// must go to stdout until then.
func realMain() int {
	// Get the command line args.
	args := os.Args[1:]

	// Initialize the application configuration.
	app.SetCfg()

	// Setup the args, Commands, and Help info.
	cli := &cli.CLI{
		Name:     app.Name,
		Version:  Version,
		Commands: Commands,
		Args:     args,
		HelpFunc: cli.BasicHelpFunc(app.Name),
	}

	// Run the passed command, recieve back a message and error object.
	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	// Return the exitcode.
	return exitCode
}

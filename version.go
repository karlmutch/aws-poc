package awstest

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/mgutz/logxi/v1"
)

// The following variables are available at link time for the build
// system to inject version and build informational values into

// BuildTime is a string with a formatted date time intended to be read by a human
var BuildTime string

// Version is the tagged version if it exists from git
var Version string

// Branch is the git branch from the build system
var Branch string

var verFlag = flag.Bool("version", false, "Print the version strings for the command, then exit")

// HandleVersion will check for the version flag string and will use it as an indication that printing
// the version is all that the user required then will bail if it is set
//
func HandleVersion() {
	log.Info(fmt.Sprintf("%s %s from %s compiled at %s using %s", os.Args[0], Version, Branch, BuildTime, runtime.Version()))
	if *verFlag {
		os.Exit(0)
	}
}

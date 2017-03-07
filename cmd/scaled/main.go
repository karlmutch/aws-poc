package main

// This file implements a simple login example for AWS

import (
	"awstest"
	"flag"
	"fmt"

	"github.com/mgutz/logxi/v1"
)

func main() {

	flag.Parse()

	// Always output the version string, however if the version flag was set by the user then this function
	// will cause the program to exit
	//
	awstest.HandleVersion()

	rdsID, cleanup, err := awstest.CreateDB("db.t2.micro", 10, "darkcycle")
	if err != nil {
		log.Warn(fmt.Sprintf("%s", err.Error()), "error", err)
	}

	// Calling this will destroy any resources allocated by the createDB including
	// security groups etc.  Calling this is optional as one might wish to keep
	// the DB around after making this request beyond the life of this execution.
	//
	if err = cleanup(); err != nil {
		log.Warn(fmt.Sprintf("%s", err.Error()), "error", err)
	}

	log.Info(fmt.Sprintf("%v", rdsID))
}

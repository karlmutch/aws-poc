package main

// This file implements a simple login example for AWS

import (
	"awstest"
	"flag"
	"fmt"
	"github.com/mgutz/logxi/v1"
	"os"
)

func main() {
	flag.Parse()

	// Always output the version string, however if the version flag was set by the user then this function
	// will cause the program to exit
	//
	awstest.HandleVersion()

	log.Info(fmt.Sprintf("%v", os.Args))
}

package main

// This file implements a simple login example for AWS

import (
	"flags"
	"fmt"
	"github.com/mgutz/logxi/v1"
	"os"
)

func main() {
	flags.Parse()

	log.Info(fmt.Sprintf("%v", os.Args))
}

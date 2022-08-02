package main

import (
	"cl-tool/changelog"
	"cl-tool/entry"
	"cl-tool/release"
	"flag"
	"fmt"
	"os"
)

const _VERSION = "1.0.3"

func main() {
	var printHelp bool

	mainCmd := flag.NewFlagSet("cl-tool "+_VERSION, flag.ExitOnError)
	mainCmd.Usage = func() { Usage(mainCmd) }
	mainCmd.BoolVar(&printHelp, "help", false, "print the help for the current command")

	if len(os.Args) <= 1 {
		mainCmd.Usage()
		os.Exit(1)
	}

	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		fmt.Printf("failed to run command: %v\n", err)
		os.Exit(1)
	}

	if printHelp {
		mainCmd.Usage()
		os.Exit(0)
	}

	cmd := os.Args[1]
	var topError error

	switch cmd {
	case changelog.CMD:
		topError = changelog.RunCmd(os.Args[2:])
	case entry.CMD:
		topError = entry.RunCmd(os.Args[2:])
	case release.CMD:
		topError = release.RunCmd(os.Args[2:])
	default:
		mainCmd.Usage()
		os.Exit(1)
	}

	if topError != nil {
		fmt.Printf("%v\n", topError)
		os.Exit(1)
	}
}

func Usage(c *flag.FlagSet) {
	fmt.Fprintf(c.Output(), "Usage of %s:\n\n", c.Name())

	fmt.Fprintf(c.Output(), "Commands:\n")
	fmt.Fprintf(c.Output(), "  generate    Output the current state of the Changelog to stdout\n")
	fmt.Fprintf(c.Output(), "  new         Create a new Changelog entry in $EDITOR, populates the 'Unreleased' folder\n")
	fmt.Fprintf(c.Output(), "  release     Rename the '.changelog/Unreleased' folder with the current date and version number\n")

	fmt.Fprintf(c.Output(), "\nOptions:\n")

	c.PrintDefaults()
}

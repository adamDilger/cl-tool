package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type BaseParams struct {
	Help bool
	Root string
}

func (n *BaseParams) NewSubCmd(name string) *flag.FlagSet {
	c := flag.NewFlagSet(name, flag.ExitOnError)
	c.Usage = func() {
		fmt.Fprintf(c.Output(), "Usage of cl-tool %s:\n", name)
		c.PrintDefaults()
	}

	c.BoolVar(&n.Help, "help", false, "print the help for this command")
	c.StringVar(&n.Root, "path", ".", "Path to the repository containing your .changelog")

	return c
}

func (n *BaseParams) ValidateWorkingDirectory() {
	// .changelog folder does not exist
	clFilePath := filepath.Join(n.Root, ".changelog")
	if _, err := os.Stat(clFilePath); os.IsNotExist(err) {
		fmt.Printf("'%s' folder does not exist. Check cwd or create directory before running.\n", clFilePath)
		os.Exit(1)
	}
}

package main

import (
	"cl-tool/changelog"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Command string

var (
	path string

	generateCmd = flag.NewFlagSet("generate", flag.ExitOnError)
	newCmd      = flag.NewFlagSet("new", flag.ExitOnError)
	releaseCmd  = flag.NewFlagSet("release", flag.ExitOnError)
)

var subcommands = map[string]*flag.FlagSet{
	generateCmd.Name(): generateCmd,
	newCmd.Name():      newCmd,
	releaseCmd.Name():  releaseCmd,
}

func setupCommonFlags() {
	for _, fs := range subcommands {
		fs.StringVar(&path, "path", ".", "Path to the repository containing your .changelog")
	}
}

func main() {
	setupCommonFlags()

	if len(os.Args) < 2 {
		fmt.Println("No subcommand passed. Must pass either [generate, new, release] as the first argument")
		os.Exit(1)
	}

	cmd := subcommands[os.Args[1]]
	if cmd == nil {
		fmt.Printf("Unknown subcommand '%s'. Must be one of [generate, new, release]\n", os.Args[1])
		os.Exit(1)
	}

	cmd.Parse(os.Args[2:])

	// .changelog folder does not exist
	clFilePath := filepath.Join(path, ".changelog")
	if _, err := os.Stat(clFilePath); os.IsNotExist(err) {
		fmt.Printf("'%s' folder does not exist. Check cwd or create directory before running.\n", clFilePath)
		os.Exit(1)
	}

	var err error

	switch cmd.Name() {
	case generateCmd.Name():
		var c *changelog.Changelog
		c, err = changelog.NewChangelog(path)
		if err != nil {
			break
		}

		c.Render(os.Stdout)
	case newCmd.Name():
		err = CreateChangelogEntry(path)
	case releaseCmd.Name():
		err = CreateRelease(path)
	default:
		flag.PrintDefaults()
	}

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

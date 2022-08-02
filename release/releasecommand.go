package release

import (
	"cl-tool/commands"
	"fmt"
)

const CMD = "release"

type releaseParams struct {
	commands.BaseParams
	Version string
}

func RunCmd(args []string) error {
	var params releaseParams
	var cmd = params.NewSubCmd(CMD)
	cmd.StringVar(&params.Version, "version", "", "set this to skip the version number prompt e.g. --version 1.0.0")

	if err := cmd.Parse(args); err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	if params.Help {
		cmd.Usage()
		return nil
	}

	params.ValidateWorkingDirectory()

	return CreateRelease(params.Root, params.Version)
}

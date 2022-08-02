package entry

import (
	"cl-tool/commands"
	"fmt"
)

const CMD = "new"

type newParams struct{ commands.BaseParams }

func RunCmd(args []string) error {
	var params newParams
	var cmd = params.NewSubCmd(CMD)

	if err := cmd.Parse(args); err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	if params.Help {
		cmd.Usage()
		return nil
	}

	params.ValidateWorkingDirectory()

	return params.NewEntry()
}

func (n *newParams) NewEntry() error {
	return CreateChangelogEntry(n.Root)
}

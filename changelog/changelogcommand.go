package changelog

import (
	"cl-tool/commands"
	"flag"
	"fmt"
	"os"
)

const CMD = "generate"

type generateParams struct{ commands.BaseParams }

func newGenerateCmd(params *generateParams) *flag.FlagSet {
	return params.NewSubCmd(CMD)
}

func RunCmd(args []string) error {
	var params generateParams
	var cmd = newGenerateCmd(&params)

	if err := cmd.Parse(args); err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	if params.Help {
		cmd.Usage()
		return nil
	}

	params.ValidateWorkingDirectory()

	return params.Generate()
}

func (n *generateParams) Generate() error {
	c, err := NewChangelog(n.Root)
	if err != nil {
		return err
	}

	c.Render(os.Stdout)
	return nil
}

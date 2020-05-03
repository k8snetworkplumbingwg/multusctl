package commands

import (
	"github.com/tliron/multusctl/version"
)

func init() {
	rootCommand.AddCommand(version.NewCommand(toolName))
}

package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/multusctl/multusctl/commands"
)

func main() {
	commands.Execute()
	atexit.Exit(0)
}

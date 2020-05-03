package commands

import (
	"github.com/op/go-logging"
)

const toolName = "multusctl"

var log = logging.MustGetLogger(toolName)

var installationNamespace string

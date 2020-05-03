package commands

import (
	"github.com/op/go-logging"
	puccinicommon "github.com/tliron/puccini/common"
	turandotcommon "github.com/tliron/turandot/common"
)

const toolName = "multusctl"

var log = logging.MustGetLogger(toolName)

var installationNamespace string

func GetNamespace(namespace string) string {
	if namespace == "" {
		if namespace, _ = turandotcommon.GetConfiguredNamespace(kubeconfigPath); namespace == "" {
			puccinicommon.Fail("could not discover namespace and \"--namespace\" not provided")
		}
	}
	return namespace
}

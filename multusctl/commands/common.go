package commands

import (
	"github.com/op/go-logging"
	"github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
)

const toolName = "multusctl"

var log = logging.MustGetLogger(toolName)

var installationNamespace string

func GetNamespace(namespace string) string {
	if namespace == "" {
		if namespace, _ = kubernetes.GetConfiguredNamespace(kubeconfigPath, context); namespace == "" {
			util.Fail("could not discover namespace and \"--namespace\" not provided")
		}
	}
	return namespace
}

package commands

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
)

const toolName = "multusctl"

var log = commonlog.GetLogger(toolName)

var installationNamespace string
var format string

func GetNamespace(namespace string) string {
	if namespace == "" {
		if namespace, _ = kubernetes.GetConfiguredNamespace(kubeconfigPath, context); namespace == "" {
			util.Fail("could not discover namespace and \"--namespace\" not provided")
		}
	}
	return namespace
}

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/multusctl/client"
	puccinicommon "github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
)

var getNamespace string

func init() {
	rootCommand.AddCommand(getCommand)
	getCommand.PersistentFlags().StringVarP(&getNamespace, "namespace", "n", "", "namespace")
}

var getCommand = &cobra.Command{
	Use:   "get [NAME]",
	Short: "Get the configuration of a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(getNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		puccinicommon.FailOnError(err)
		networkAttachmentDefinition, err := client.Get(args[0])
		puccinicommon.FailOnError(err)
		data, err := format.DecodeJSON(networkAttachmentDefinition.Spec.Config)
		puccinicommon.FailOnError(err)
		config, err := format.EncodeYAML(data, "  ", false)
		puccinicommon.FailOnError(err)
		fmt.Fprint(terminal.Stdout, config)
	},
}

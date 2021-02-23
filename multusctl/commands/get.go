package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/multusctl/client"
)

var getNamespace string

func init() {
	rootCommand.AddCommand(getCommand)
	getCommand.Flags().StringVarP(&getNamespace, "namespace", "n", "", "namespace")
}

var getCommand = &cobra.Command{
	Use:   "get [NAME]",
	Short: "Get the configuration of a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(getNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		util.FailOnError(err)
		networkAttachmentDefinition, err := client.GetNetworkAttachmentDefinition(args[0])
		util.FailOnError(err)
		data, _, err := ard.DecodeJSON(networkAttachmentDefinition.Spec.Config, false)
		util.FailOnError(err)
		config, err := format.EncodeYAML(data, "  ", false)
		util.FailOnError(err)
		fmt.Fprint(terminal.Stdout, config)
	},
}

package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
	"github.com/tliron/multusctl/client"
)

var deleteNamespace string

func init() {
	rootCommand.AddCommand(deleteCommand)
	deleteCommand.Flags().StringVarP(&deleteNamespace, "namespace", "n", "", "namespace")
}

var deleteCommand = &cobra.Command{
	Use:   "delete [NAME]",
	Short: "Delete a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(deleteNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		util.FailOnError(err)
		err = client.DeleteNetworkAttachmentDefinition(args[0])
		util.FailOnError(err)
	},
}

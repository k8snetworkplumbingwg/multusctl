package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/multusctl/client"
	puccinicommon "github.com/tliron/puccini/common"
)

var deleteNamespace string

func init() {
	rootCommand.AddCommand(deleteCommand)
	deleteCommand.PersistentFlags().StringVarP(&deleteNamespace, "namespace", "n", "", "namespace")
}

var deleteCommand = &cobra.Command{
	Use:   "delete [NAME]",
	Short: "Delete a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(deleteNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		puccinicommon.FailOnError(err)
		err = client.Delete(args[0])
		puccinicommon.FailOnError(err)
	},
}

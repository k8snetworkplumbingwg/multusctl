package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/multusctl/client"
	"github.com/tliron/puccini/common"
)

func init() {
	rootCommand.AddCommand(uninstallCommand)
	uninstallCommand.PersistentFlags().StringVarP(&installationNamespace, "namespace", "n", "kube-system", "namespace")
}

var uninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Multus CNI",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClient(masterUrl, kubeconfigPath, installationNamespace)
		common.FailOnError(err)
		client.Uninstall()
	},
}

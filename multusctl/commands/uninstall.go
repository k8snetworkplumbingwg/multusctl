package commands

import (
	"github.com/k8snetworkplumbingwg/multusctl/client"
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	rootCommand.AddCommand(uninstallCommand)
	uninstallCommand.Flags().StringVarP(&installationNamespace, "namespace", "n", "kube-system", "namespace")
}

var uninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Multus CNI",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClient(masterUrl, kubeconfigPath, context, installationNamespace)
		util.FailOnError(err)
		client.Uninstall()
	},
}

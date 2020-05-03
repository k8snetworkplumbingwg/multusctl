package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/multusctl/client"
	puccinicommon "github.com/tliron/puccini/common"
)

var registry string
var wait bool

func init() {
	rootCommand.AddCommand(installCommand)
	installCommand.PersistentFlags().StringVarP(&installationNamespace, "namespace", "n", "kube-system", "namespace")
	installCommand.PersistentFlags().StringVarP(&registry, "registry", "r", "docker.io", "registry URL (use special value \"internal\" to discover internally deployed registry)")
	installCommand.PersistentFlags().BoolVarP(&wait, "wait", "w", false, "wait for installation to succeed")
}

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Install Multus CNI",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClient(masterUrl, kubeconfigPath, installationNamespace)
		puccinicommon.FailOnError(err)
		err = client.Install(registry, wait)
		puccinicommon.FailOnError(err)
	},
}

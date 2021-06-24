package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
	"github.com/tliron/multusctl/client"
)

var runtime string
var registry string
var wait bool

func init() {
	rootCommand.AddCommand(installCommand)
	installCommand.Flags().StringVarP(&installationNamespace, "namespace", "n", "kube-system", "namespace")
	installCommand.Flags().StringVarP(&runtime, "runtime", "t", "crio", "container runtime (\"crio\" or \"default\")")
	installCommand.Flags().StringVarP(&registry, "registry", "r", "ghcr.io", "registry URL (use special value \"internal\" to discover internally deployed registry)")
	installCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for installation to succeed")
}

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Install Multus CNI",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClient(masterUrl, kubeconfigPath, installationNamespace)
		util.FailOnError(err)
		err = client.Install(runtime, registry, wait)
		util.FailOnError(err)
	},
}

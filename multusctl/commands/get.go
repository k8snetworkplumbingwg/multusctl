package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/multusctl/client"
)

var getNamespace string
var strict bool
var pretty bool

func init() {
	rootCommand.AddCommand(getCommand)
	getCommand.Flags().StringVarP(&getNamespace, "namespace", "n", "", "namespace")
	getCommand.Flags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", \"cjson\", \"xml\", or \"cbor\")")
	rootCommand.PersistentFlags().BoolVarP(&strict, "strict", "y", false, "strict output (for \"YAML\" format only)")
	getCommand.Flags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")
}

var getCommand = &cobra.Command{
	Use:   "get [NAME]",
	Short: "Get the configuration of a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(getNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, context, namespace)
		util.FailOnError(err)
		networkAttachmentDefinition, err := client.GetNetworkAttachmentDefinition(args[0])
		util.FailOnError(err)
		data, _, err := ard.DecodeJSON(networkAttachmentDefinition.Spec.Config, false)
		util.FailOnError(err)
		data, _ = ard.MapsToStringMaps(data)
		err = formatpkg.Print(data, format, terminal.Stdout, strict, pretty)
		util.FailOnError(err)
	},
}

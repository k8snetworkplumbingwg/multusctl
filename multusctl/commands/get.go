package commands

import (
	"os"

	"github.com/k8snetworkplumbingwg/multusctl/client"
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
)

var getNamespace string
var strict bool
var pretty bool

func init() {
	rootCommand.AddCommand(getCommand)
	getCommand.Flags().StringVarP(&getNamespace, "namespace", "n", "", "namespace")
	getCommand.Flags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", \"cjson\", \"xml\", \"cbor\", \"messagepack\", or \"go\")")
	rootCommand.PersistentFlags().BoolVarP(&strict, "strict", "y", false, "strict output (for \"yaml\" format only)")
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
		data, _ = ard.NormalizeStringMaps(data)
		err = transcribe.Print(data, format, os.Stdout, strict, pretty)
		util.FailOnError(err)
	},
}

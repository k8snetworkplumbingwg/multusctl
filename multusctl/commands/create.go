package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
	"github.com/tliron/multusctl/client"
)

var createNamespace string
var configFile string

func init() {
	rootCommand.AddCommand(createCommand)
	createCommand.Flags().StringVarP(&createNamespace, "namespace", "n", "", "namespace")
	createCommand.Flags().StringVarP(&configFile, "file", "f", "", "path to config file (YAML or JSON)")
}

var createCommand = &cobra.Command{
	Use:   "create [NAME]",
	Short: "Create a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(createNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		util.FailOnError(err)

		var config string

		if configFile != "" {
			urlContext := urlpkg.NewContext()
			defer urlContext.Release()

			url := urlpkg.NewFileURL(configFile, urlContext)
			config, err = urlpkg.ReadString(url)

			switch url.Format() {
			case "json":
				err = format.ValidateJSON(config)
				util.FailOnError(err)

			case "yaml":
				data, _, err := ard.DecodeYAML(config, false)
				util.FailOnError(err)
				data, _ = ard.MapsToStringMaps(data)
				config, err = format.EncodeJSON(data, "  ")
				util.FailOnError(err)
			}
		} else {
			util.Fail("must provide \"--file\" or TODO")
		}

		_, err = client.CreateNetworkAttachmentDefinition(args[0], config)
		util.FailOnError(err)
	},
}

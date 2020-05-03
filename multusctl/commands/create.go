package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/multusctl/client"
	"github.com/tliron/puccini/ard"
	puccinicommon "github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/format"
	urlpkg "github.com/tliron/puccini/url"
)

var createNamespace string
var configFile string

func init() {
	rootCommand.AddCommand(createCommand)
	createCommand.PersistentFlags().StringVarP(&createNamespace, "namespace", "n", "", "namespace")
	createCommand.PersistentFlags().StringVarP(&configFile, "file", "f", "", "path to config file (YAML or JSON)")
}

var createCommand = &cobra.Command{
	Use:   "create [NAME]",
	Short: "Create a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(createNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		puccinicommon.FailOnError(err)

		var config string

		if configFile != "" {
			url := urlpkg.NewFileURL(configFile)
			config, err = urlpkg.ReadToString(url)

			switch url.Format() {
			case "json":
				err = format.ValidateJSON(config)
				puccinicommon.FailOnError(err)
			case "yaml":
				data, err := format.DecodeYAML(config)
				puccinicommon.FailOnError(err)
				data, _ = ard.ToStringMaps(data)
				config, err = format.EncodeJSON(data, "  ")
				puccinicommon.FailOnError(err)
			}
		} else {
			puccinicommon.Fail("must provide \"--file\" or TODO")
		}

		_, err = client.Create(args[0], config)
		puccinicommon.FailOnError(err)
	},
}

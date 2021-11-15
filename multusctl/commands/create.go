package commands

import (
	"github.com/k8snetworkplumbingwg/multusctl/client"
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	formatpkg "github.com/tliron/kutil/format"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
)

var createNamespace string
var configUrl string

func init() {
	rootCommand.AddCommand(createCommand)
	createCommand.Flags().StringVarP(&createNamespace, "namespace", "n", "", "namespace")
	createCommand.Flags().StringVarP(&configUrl, "url", "u", "", "URL or path to config file (defaults to stdin)")
	createCommand.Flags().StringVarP(&format, "format", "f", "", "force input format (\"yaml\" or \"json\", defaults to URL extension)")
}

var createCommand = &cobra.Command{
	Use:   "create [NAME]",
	Short: "Create a network attachment definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(createNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, context, namespace)
		util.FailOnError(err)

		var url urlpkg.URL

		if configUrl != "" {
			urlContext := urlpkg.NewContext()
			util.OnExit(func() {
				if err := urlContext.Release(); err != nil {
					log.Errorf("%s", err.Error())
				}
			})

			url, err = urlpkg.NewValidURL(configUrl, nil, urlContext)
			util.FailOnError(err)
			if format == "" {
				format = url.Format()
			}
		} else {
			if format == "" {
				format = "yaml"
			}
			url, err = urlpkg.ReadToInternalURLFromStdin(format)
			util.FailOnError(err)
		}

		var config string
		config, err = urlpkg.ReadString(url)
		util.FailOnError(err)

		switch format {
		case "json":
			err = formatpkg.ValidateJSON(config)
			util.FailOnError(err)

		case "yaml":
			data, _, err := ard.DecodeYAML(config, false)
			util.FailOnError(err)
			data, _ = ard.MapsToStringMaps(data)
			config, err = formatpkg.EncodeJSON(data, "  ")
			util.FailOnError(err)

		default:
			util.Failf("unsupported format: %q", format)
		}

		_, err = client.CreateNetworkAttachmentDefinition(args[0], config)
		util.FailOnError(err)
	},
}

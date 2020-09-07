package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/multusctl/client"
)

var listNamespace string
var bare bool

func init() {
	rootCommand.AddCommand(listCommand)
	listCommand.Flags().StringVarP(&listNamespace, "namespace", "n", "", "namespace")
	listCommand.Flags().BoolVarP(&bare, "bare", "b", false, "list bare names (not as a table)")
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List network attachment definitions",
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(listNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		util.FailOnError(err)
		networkAttachmentDefintions, err := client.ListNetworkAttachmentDefinitions()
		util.FailOnError(err)

		if bare {
			for _, networkAttachmentDefintion := range networkAttachmentDefintions.Items {
				fmt.Fprintln(terminal.Stdout, networkAttachmentDefintion.Name)
			}
		} else {
			table := terminal.NewTable(maxWidth, "Name", "Pods")
			for _, networkAttachmentDefintion := range networkAttachmentDefintions.Items {
				pods, err := client.ListPodsForNetworkAttachmentDefinition(networkAttachmentDefintion.Name)
				util.FailOnError(err)
				podNames := ""
				for _, pod := range pods {
					podNames += pod.Name + "\n"
				}
				table.Add(networkAttachmentDefintion.Name, podNames)
			}
			table.Print()
		}
	},
}

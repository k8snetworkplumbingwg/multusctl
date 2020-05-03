package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/multusctl/client"
	puccinicommon "github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/terminal"
	turandotcommon "github.com/tliron/turandot/common"
)

var listNamespace string
var bare bool

func init() {
	rootCommand.AddCommand(listCommand)
	listCommand.PersistentFlags().StringVarP(&listNamespace, "namespace", "n", "", "namespace")
	listCommand.PersistentFlags().BoolVarP(&bare, "bare", "b", false, "list bare names (not as a table)")
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List network attachment definitions",
	Run: func(cmd *cobra.Command, args []string) {
		namespace := GetNamespace(listNamespace)
		client, err := client.NewClient(masterUrl, kubeconfigPath, namespace)
		puccinicommon.FailOnError(err)
		networkAttachmentDefintions, err := client.List()
		puccinicommon.FailOnError(err)

		if bare {
			for _, networkAttachmentDefintion := range networkAttachmentDefintions.Items {
				fmt.Fprintln(terminal.Stdout, networkAttachmentDefintion.Name)
			}
		} else {
			table := turandotcommon.NewTable("Name", "Pods")
			for _, networkAttachmentDefintion := range networkAttachmentDefintions.Items {
				pods, err := client.ListPods(networkAttachmentDefintion.Name)
				puccinicommon.FailOnError(err)
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

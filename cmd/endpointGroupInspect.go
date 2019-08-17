package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/greenled/portainer-stack-utils/common"
	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// endpointGroupInspectCmd represents the endpoint group inspect command
var endpointGroupInspectCmd = &cobra.Command{
	Use:     "inspect <name>",
	Short:   "Inspect an endpoint group",
	Example: "  psu endpoint group inspect production",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpointGroupName := args[0]
		endpointGroup, err := common.GetEndpointGroupByName(endpointGroupName)
		if err == nil {
			// The endpoint group exists
			switch viper.GetString("endpoint.group.inspect.format") {
			case "table":
				// Print endpoint group in a table format
				writer, err := common.NewTabWriter([]string{
					"ID",
					"NAME",
					"DESCRIPTION",
				})
				common.CheckError(err)
				_, err = fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%s",
					endpointGroup.ID,
					endpointGroup.Name,
					endpointGroup.Description,
				))
				common.CheckError(err)
				err = writer.Flush()
				common.CheckError(err)
			case "json":
				// Print endpoint group in a json format
				endpointJsonBytes, err := json.Marshal(endpointGroup)
				common.CheckError(err)
				fmt.Println(string(endpointJsonBytes))
			default:
				// Print endpoint group in a custom format
				template, err := template.New("endpointGroupTpl").Parse(viper.GetString("endpoint.group.inspect.format"))
				common.CheckError(err)
				err = template.Execute(os.Stdout, endpointGroup)
				common.CheckError(err)
				fmt.Println()
			}
		} else if err == common.ErrEndpointGroupNotFound {
			// The endpoint group does not exist
			logrus.WithFields(logrus.Fields{
				"endpointGroup": endpointGroupName,
			}).Fatal("Endpoint group not found")
		} else {
			// Something else happened
			common.CheckError(err)
		}
	},
}

func init() {
	endpointGroupCmd.AddCommand(endpointGroupInspectCmd)

	endpointGroupInspectCmd.Flags().String("format", "table", `Output format. Can be "table", "json" or a Go template.`)
	viper.BindPFlag("endpoint.group.inspect.format", endpointGroupInspectCmd.Flags().Lookup("format"))

	endpointGroupInspectCmd.SetUsageTemplate(endpointGroupInspectCmd.UsageTemplate() + common.GetFormatHelp(portainer.EndpointGroup{}))
}

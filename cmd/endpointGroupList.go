package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// endpointGroupListCmd represents the endpoint group list command
var endpointGroupListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List endpoint groups",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := common.GetClient()
		common.CheckError(err)

		logrus.Debug("Getting endpoint groups")
		endpointGroups, err := client.GetEndpointGroups()
		common.CheckError(err)

		if viper.GetString("endpoint.group.list.format") != "" {
			// Print endpoint group fields formatted
			template, templateParsingErr := template.New("endpointGroupTpl").Parse(viper.GetString("endpoint.group.list.format"))
			common.CheckError(templateParsingErr)
			for _, g := range endpointGroups {
				templateExecutionErr := template.Execute(os.Stdout, g)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		} else {
			// Print all endpoint group fields as a table
			writer, err := common.NewTabWriter([]string{
				"ID",
				"NAME",
				"DESCRIPTION",
			})
			common.CheckError(err)
			for _, g := range endpointGroups {
				_, err := fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%s",
					g.Id,
					g.Name,
					g.Description,
				))
				common.CheckError(err)
			}
			flushErr := writer.Flush()
			common.CheckError(flushErr)
		}
	},
}

func init() {
	endpointGroupCmd.AddCommand(endpointGroupListCmd)

	endpointGroupListCmd.Flags().String("format", "", "Format output using a Go template.")
	viper.BindPFlag("endpoint.group.list.format", endpointGroupListCmd.Flags().Lookup("format"))
}

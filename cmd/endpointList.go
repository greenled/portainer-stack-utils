package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	portainer "github.com/portainer/portainer/api"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// endpointListCmd represents the list command
var endpointListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List endpoints",
	Aliases: []string{"ls"},
	Example: `  Print endpoints in a table format:
  psu endpoint ls

  Print whether an endpoint is a Swarm cluster or not:
  psu endpoint ls --format "{{ .Name }} ({{ .ID }}): {{ range .Snapshots }}{{ if .Swarm }}yes{{ else }}no{{ end }}{{ end }}"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := common.GetClient()
		common.CheckError(err)

		logrus.Debug("Getting endpoints")
		endpoints, err := client.GetEndpoints()
		common.CheckError(err)

		switch viper.GetString("endpoint.list.format") {
		case "table":
			// Print endpoints in a table format
			writer, err := common.NewTabWriter([]string{
				"ID",
				"NAME",
				"TYPE",
				"URL",
				"PUBLIC URL",
				"GROUP ID",
			})
			common.CheckError(err)
			for _, e := range endpoints {
				var endpointType string
				if e.Type == portainer.DockerEnvironment {
					endpointType = "docker"
				} else if e.Type == portainer.AgentOnDockerEnvironment {
					endpointType = "agent on docker"
				} else if e.Type == portainer.AzureEnvironment {
					endpointType = "azure"
				} else if e.Type == portainer.EdgeAgentEnvironment {
					endpointType = "edge agent"
				}
				_, err := fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%v\t%s\t%s\t%v",
					e.ID,
					e.Name,
					endpointType,
					e.URL,
					e.PublicURL,
					e.GroupID,
				))
				common.CheckError(err)
			}
			flushErr := writer.Flush()
			common.CheckError(flushErr)
		case "json":
			// Print endpoints in a json format
			statusJsonBytes, err := json.Marshal(endpoints)
			common.CheckError(err)
			fmt.Println(string(statusJsonBytes))
		default:
			// Print endpoints in a custom format
			template, templateParsingErr := template.New("endpointTpl").Parse(viper.GetString("endpoint.list.format"))
			common.CheckError(templateParsingErr)
			for _, e := range endpoints {
				templateExecutionErr := template.Execute(os.Stdout, e)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		}
	},
}

func init() {
	endpointCmd.AddCommand(endpointListCmd)

	endpointListCmd.Flags().String("format", "table", `Output format. Can be "table", "json" or a Go template.`)
	viper.BindPFlag("endpoint.list.format", endpointListCmd.Flags().Lookup("format"))

	endpointListCmd.SetUsageTemplate(endpointListCmd.UsageTemplate() + common.GetFormatHelp(portainer.Endpoint{}))
}

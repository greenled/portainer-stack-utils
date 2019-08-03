package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/greenled/portainer-stack-utils/util"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stackListCmd represents the remove command
var stackListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List stacks",
	Aliases: []string{"ls"},
	Example: "psu stack list --endpoint 1",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := common.GetClient()
		util.CheckError(err)

		stacks, err := client.GetStacks(viper.GetString("stack.list.swarm"), viper.GetUint32("stack.list.endpoint"))
		util.CheckError(err)

		if viper.GetBool("stack.list.quiet") {
			// Print only stack names
			for _, s := range stacks {
				_, err := fmt.Println(s.Name)
				util.CheckError(err)
			}
		} else if viper.GetString("stack.list.format") != "" {
			// Print stack fields formatted
			template, templateParsingErr := template.New("stackTpl").Parse(viper.GetString("stack.list.format"))
			util.CheckError(templateParsingErr)
			for _, s := range stacks {
				templateExecutionErr := template.Execute(os.Stdout, s)
				util.CheckError(templateExecutionErr)
				fmt.Println()
			}
		} else {
			// Print all stack fields as a table
			writer, err := util.NewTabWriter([]string{
				"STACK ID",
				"NAME",
				"TYPE",
				"ENTRY POINT",
				"PROJECT PATH",
				"ENDPOINT ID",
				"SWARM ID",
			})
			util.CheckError(err)
			for _, s := range stacks {
				_, err := fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%v\t%s\t%s\t%v\t%s",
					s.Id,
					s.Name,
					s.GetTranslatedStackType(),
					s.EntryPoint,
					s.ProjectPath,
					s.EndpointID,
					s.SwarmID,
				))
				util.CheckError(err)
			}
			flushErr := writer.Flush()
			util.CheckError(flushErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackListCmd)

	stackListCmd.Flags().String("swarm", "", "filter by swarm ID")
	stackListCmd.Flags().String("endpoint", "", "filter by endpoint ID")
	stackListCmd.Flags().BoolP("quiet", "q", false, "only display stack names")
	stackListCmd.Flags().String("format", "", "format output using a Go template")
	viper.BindPFlag("stack.list.swarm", stackListCmd.Flags().Lookup("swarm"))
	viper.BindPFlag("stack.list.endpoint", stackListCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.list.quiet", stackListCmd.Flags().Lookup("quiet"))
	viper.BindPFlag("stack.list.format", stackListCmd.Flags().Lookup("format"))
}

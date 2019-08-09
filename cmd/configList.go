package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/greenled/portainer-stack-utils/common"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// configListCmd represents the list command
var configListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List configs",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get alphabetically ordered list of config keys
		keys := viper.AllKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		// Create config objects
		var configs []config
		for _, key := range keys {
			envvar := strings.Replace(key, "-", "_", -1)
			envvar = strings.Replace(envvar, ".", "_", -1)
			envvar = strings.ToUpper(envvar)
			envvar = "PSU_" + envvar
			configs = append(configs, config{
				Key:                 key,
				EnvironmentVariable: envvar,
				CurrentValue:        viper.Get(key),
			})
		}

		if viper.GetString("config.list.format") != "" {
			// Print configs with a custom format
			template, templateParsingErr := template.New("configTpl").Parse(viper.GetString("config.list.format"))
			common.CheckError(templateParsingErr)
			for _, c := range configs {
				templateExecutionErr := template.Execute(os.Stdout, c)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		} else {
			// Print configs with a table format
			writer, err := common.NewTabWriter([]string{
				"KEY",
				"ENV VAR",
				"CURRENT VALUE",
			})
			common.CheckError(err)
			for _, c := range configs {
				_, err := fmt.Fprintln(writer, fmt.Sprintf(
					"%s\t%s\t%v",
					c.Key,
					c.EnvironmentVariable,
					c.CurrentValue,
				))
				common.CheckError(err)
			}
			flushErr := writer.Flush()
			common.CheckError(flushErr)
		}
	},
}

func init() {
	configCmd.AddCommand(configListCmd)

	configListCmd.Flags().String("format", "", "Format output using a Go template.")
	viper.BindPFlag("config.list.format", configListCmd.Flags().Lookup("format"))
}

type config struct {
	Key                 string
	EnvironmentVariable string
	CurrentValue        interface{}
}

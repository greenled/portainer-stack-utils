package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/greenled/portainer-stack-utils/common"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// settingListCmd represents the setting list command
var settingListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List settings",
	Aliases: []string{"ls"},
	Example: `  Print settings in a table format:
  psu setting ls

  Print available setting keys:
  psu setting ls --format "{{ .Key }}"

  Print settings in a yaml|properties format:
  psu setting ls --format "{{ .Key }}:{{ if .CurrentValue }} {{ .CurrentValue }}{{ end }}"

  Print available environment variables:
  psu setting ls --format "{{ .EnvironmentVariable }}"  

  Print settings in a dotenv format:
  psu setting ls --format "{{ .EnvironmentVariable }}={{ if .CurrentValue }}{{ .CurrentValue }}{{ end }}"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get alphabetically ordered list of setting keys
		keys := viper.AllKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		// Create setting objects
		var settings []setting
		for _, key := range keys {
			envvar := strings.Replace(key, "-", "_", -1)
			envvar = strings.Replace(envvar, ".", "_", -1)
			envvar = strings.ToUpper(envvar)
			envvar = "PSU_" + envvar
			settings = append(settings, setting{
				Key:                 key,
				EnvironmentVariable: envvar,
				CurrentValue:        viper.Get(key),
			})
		}

		switch viper.GetString("setting.list.format") {
		case "table":
			// Print settings in a table format
			writer, err := common.NewTabWriter([]string{
				"KEY",
				"ENV VAR",
				"CURRENT VALUE",
			})
			common.CheckError(err)
			for _, c := range settings {
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
		case "json":
			// Print settings in a json format
			statusJSONBytes, err := json.Marshal(settings)
			common.CheckError(err)
			fmt.Println(string(statusJSONBytes))
		default:
			// Print settings in a custom format
			template, templateParsingErr := template.New("settingTpl").Parse(viper.GetString("setting.list.format"))
			common.CheckError(templateParsingErr)
			for _, c := range settings {
				templateExecutionErr := template.Execute(os.Stdout, c)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		}
	},
}

func init() {
	settingCmd.AddCommand(settingListCmd)

	settingListCmd.Flags().String("format", "table", `Output format. Can be "table", "json" or a Go template.`)
	viper.BindPFlag("setting.list.format", settingListCmd.Flags().Lookup("format"))

	settingListCmd.SetUsageTemplate(settingListCmd.UsageTemplate() + common.GetFormatHelp(setting{}))
}

type setting struct {
	Key                 string
	EnvironmentVariable string
	CurrentValue        interface{}
}

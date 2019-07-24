package cmd

import (
	"fmt"
	"github.com/greenled/portainer-stack-utils/common"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config KEY [VALUE]",
	Short: "Get and set configuration options",
	Example: "psu config user admin",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if it's a valid key
		var keyExists bool
		for _, key := range viper.AllKeys() {
			if key == args[0] {
				keyExists = true
				break
			}
		}
		if !keyExists {
			log.Fatalf("Unkonwn configuration key \"%s\"", args[0])
		}

		// Create new viper
		commandViper := viper.New()

		// Set config file name
		var configFile string
		if viper.ConfigFileUsed() != "" {
			// Use config file from viper
			configFile = viper.ConfigFileUsed()
		} else {
			// Find home directory
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Use $HOME/.psu.yaml
			configFile = fmt.Sprintf("%s%s.psu.yaml", home, string(os.PathSeparator))
		}
		commandViper.SetConfigFile(configFile)

		// Read config from file
		if configReadingErr := commandViper.ReadInConfig(); configReadingErr != nil {
			common.PrintVerbose(fmt.Sprintf("Could not read configuration from \"%s\". Expect all configuration values to be unset.", configFile))
		}

		if len(args) == 1 {
			// Get config
			fmt.Println(commandViper.Get(args[0]))
		} else {
			// Set config
			commandViper.Set(args[0], args[1])

			// Make sure the config file exists
			_, fileCreationErr := os.Create(configFile)
			if fileCreationErr != nil {
				common.CheckError(fileCreationErr)
			}

			// Write te config file
			configWritingErr := commandViper.WriteConfig()
			if configWritingErr != nil {
				common.CheckError(configWritingErr)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

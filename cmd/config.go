package cmd

import (
	"fmt"
	"os"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config KEY [VALUE]",
	Short:   "Get and set configuration options",
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
			logrus.WithFields(logrus.Fields{
				"key":         args[0],
				"suggestions": "Try looking up the available configuration keys: psu config ls",
			}).Fatal("Unknown configuration key")
		}

		if len(args) == 1 {
			// Get config
			value, configGettingErr := getConfig(args[0])
			common.CheckError(configGettingErr)
			fmt.Println(value)
		} else {
			// Set config
			configSettingErr := setConfig(args[0], args[1])
			common.CheckError(configSettingErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func loadCofig() (*viper.Viper, error) {
	// Set config file name
	var configFile string
	if viper.ConfigFileUsed() != "" {
		// Use config file from viper
		configFile = viper.ConfigFileUsed()
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			return &viper.Viper{}, err
		}

		// Use $HOME/.psu.yaml
		configFile = fmt.Sprintf("%s%s.psu.yaml", home, string(os.PathSeparator))
	}
	newViper := viper.New()
	newViper.SetConfigFile(configFile)

	// Read config from file
	if configReadingErr := newViper.ReadInConfig(); configReadingErr != nil {
		logrus.WithFields(logrus.Fields{
			"file": configFile,
		}).Warn("Could not read configuration from file. Expect all configuration values to be unset.")
	}

	return newViper, nil
}

func getConfig(key string) (interface{}, error) {
	newViper, configLoadingErr := loadCofig()
	if configLoadingErr != nil {
		return nil, configLoadingErr
	}

	return newViper.Get(key), nil
}

func setConfig(key string, value string) error {
	newViper, configLoadingErr := loadCofig()
	if configLoadingErr != nil {
		return configLoadingErr
	}

	newViper.Set(key, value)

	// Make sure the config file exists
	_, fileCreationErr := os.Create(newViper.ConfigFileUsed())
	if fileCreationErr != nil {
		return fileCreationErr
	}

	// Write te config file
	configWritingErr := newViper.WriteConfig()
	if configWritingErr != nil {
		return configWritingErr
	}

	return nil
}

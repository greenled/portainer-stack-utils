package cmd

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/greenled/portainer-stack-utils/common"

	"github.com/spf13/cobra"
)

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set config",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		keyExists := common.CheckConfigKeyExists(args[0])
		if !keyExists {
			logrus.WithFields(logrus.Fields{
				"key":         args[0],
				"suggestions": "try looking up the available configuration keys: psu config ls",
			}).Fatal("unknown configuration key")
		}

		// Set config
		configSettingErr := setConfig(args[0], args[1])
		common.CheckError(configSettingErr)
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}

func setConfig(key string, value string) (err error) {
	newViper, err := common.LoadCofig()
	if err != nil {
		return
	}

	newViper.Set(key, value)

	// Make sure the config file exists
	_, err = os.Create(newViper.ConfigFileUsed())
	if err != nil {
		return
	}

	// Write te config file
	err = newViper.WriteConfig()

	return
}

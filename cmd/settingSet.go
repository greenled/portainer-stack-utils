package cmd

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/greenled/portainer-stack-utils/common"

	"github.com/spf13/cobra"
)

// settingSetCmd represents the setting set command
var settingSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set setting",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		keyExists := common.CheckSettingKeyExists(args[0])
		if !keyExists {
			logrus.WithFields(logrus.Fields{
				"key":         args[0],
				"suggestions": "try looking up the available setting keys: psu setting ls",
			}).Fatal("unknown setting key")
		}

		// Set setting
		err := setSetting(args[0], args[1])
		common.CheckError(err)
	},
}

func init() {
	settingCmd.AddCommand(settingSetCmd)
}

func setSetting(key string, value string) (err error) {
	newViper, err := common.LoadSettings()
	if err != nil {
		return
	}

	newViper.Set(key, value)

	// Make sure the setting file exists
	_, err = os.Create(newViper.ConfigFileUsed())
	if err != nil {
		return
	}

	// Write te setting file
	err = newViper.WriteConfig()

	return
}

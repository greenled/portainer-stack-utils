package cmd

import (
	"fmt"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// settingGetCmd represents the setting get command
var settingGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get setting",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyExists := common.CheckSettingKeyExists(args[0])
		if !keyExists {
			logrus.WithFields(logrus.Fields{
				"key":         args[0],
				"suggestions": "try looking up the available setting keys: psu setting ls",
			}).Fatal("unknown setting key")
		}

		// Get setting
		value, err := getSetting(args[0])
		common.CheckError(err)
		fmt.Println(value)
	},
}

func init() {
	settingCmd.AddCommand(settingGetCmd)
}

func getSetting(key string) (value interface{}, err error) {
	newViper, err := common.LoadSettings()
	if err != nil {
		return
	}
	value = newViper.Get(key)

	return
}

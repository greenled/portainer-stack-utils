package cmd

import (
	"fmt"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyExists := common.CheckConfigKeyExists(args[0])
		if !keyExists {
			logrus.WithFields(logrus.Fields{
				"key":         args[0],
				"suggestions": "Try looking up the available configuration keys: psu config ls",
			}).Fatal("Unknown configuration key")
		}

		// Get config
		value, configGettingErr := getConfig(args[0])
		common.CheckError(configGettingErr)
		fmt.Println(value)
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
}

func getConfig(key string) (value interface{}, err error) {
	newViper, err := common.LoadCofig()
	if err != nil {
		return
	}
	value = newViper.Get(key)

	return
}

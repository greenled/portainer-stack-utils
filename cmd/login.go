package cmd

import (
	"fmt"

	"github.com/greenled/portainer-stack-utils/client"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Portainer instance",
	Run: func(cmd *cobra.Command, args []string) {
		// Get auth token
		portainerClient, err := common.GetDefaultClient()
		common.CheckError(err)

		user := viper.GetString("user")
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Debug("Getting auth token")
		authToken, err := portainerClient.AuthenticateUser(client.AuthenticateUserOptions{
			Username: viper.GetString("user"),
			Password: viper.GetString("password"),
		})
		common.CheckError(err)

		if viper.GetBool("login.print") {
			fmt.Println(authToken)
		}

		// Save auth token
		settingErr := setSetting("auth-token", authToken)
		common.CheckError(settingErr)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().Bool("print", false, "Print retrieved auth token.")
	viper.BindPFlag("login.print", loginCmd.Flags().Lookup("print"))
}

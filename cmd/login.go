package cmd

import (
	"fmt"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Portainer instance",
	Run: func(cmd *cobra.Command, args []string) {
		// Get auth token
		client, err := common.GetDefaultClient()
		common.CheckError(err)

		authToken, err := client.Authenticate()
		common.CheckError(err)

		if viper.GetBool("login.print") {
			fmt.Println(authToken)
		}

		// Save auth token
		configSettingErr := setConfig("auth-token", authToken)
		common.CheckError(configSettingErr)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().Bool("print", false, "prints retrieved auth token")
	viper.BindPFlag("login.print", loginCmd.Flags().Lookup("print"))
}

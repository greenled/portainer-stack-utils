package cmd

import (
	"github.com/greenled/portainer-stack-utils/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stackRemoveCmd represents the remove command
var stackRemoveCmd = &cobra.Command{
	Use:     "remove STACK_NAME",
	Short:   "Remove a stack",
	Aliases: []string{"rm", "down"},
	Example: "psu stack rm mystack",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]
		logrus.WithFields(logrus.Fields{
			"stack": stackName,
		}).Debug("Getting stack")
		stack, err := common.GetStackByName(stackName, "", 0)

		switch err.(type) {
		case nil:
			// The stack exists
			stackId := stack.Id

			client, err := common.GetClient()
			common.CheckError(err)

			logrus.WithFields(logrus.Fields{
				"stack": stackName,
			}).Info("Removing stack")
			err = client.DeleteStack(stackId)
			common.CheckError(err)
		case *common.StackNotFoundError:
			// The stack does not exist
			logrus.WithFields(logrus.Fields{
				"stack": stackName,
			}).Debug("Stack not found")
			if viper.GetBool("stack.remove.strict") {
				logrus.WithFields(logrus.Fields{
					"stack": stackName,
				}).Fatal("Stack does not exist")
			}
		default:
			// Something else happened
			common.CheckError(err)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackRemoveCmd)

	stackRemoveCmd.Flags().Bool("strict", false, "Fail if stack does not exist.")
	viper.BindPFlag("stack.remove.strict", stackRemoveCmd.Flags().Lookup("strict"))
}

package cmd

import (
	"fmt"
	"log"

	"github.com/greenled/portainer-stack-utils/util"

	"github.com/greenled/portainer-stack-utils/common"
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
		stack, err := common.GetStackByName(stackName)

		switch err.(type) {
		case nil:
			// The stack exists
			util.PrintVerbose(fmt.Sprintf("Stack %s exists.", stackName))

			stackId := stack.Id

			util.PrintVerbose(fmt.Sprintf("Removing stack %s...", stackName))

			client, err := common.GetClient()
			common.CheckError(err)

			util.PrintVerbose("Deleting stack...")
			err = client.DeleteStack(stackId)
			common.CheckError(err)
		case *common.StackNotFoundError:
			// The stack does not exist
			util.PrintVerbose(fmt.Sprintf("Stack %s does not exist.", stackName))
			if viper.GetBool("stack.remove.strict") {
				log.Fatalln(fmt.Sprintf("Stack %s does not exist.", stackName))
			}
		default:
			// Something else happened
			common.CheckError(err)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackRemoveCmd)

	stackRemoveCmd.Flags().Bool("strict", false, "fail if stack does not exist")
	viper.BindPFlag("stack.remove.strict", stackRemoveCmd.Flags().Lookup("strict"))
}

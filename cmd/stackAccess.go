package cmd

import (
	"net/http"

	"github.com/greenled/portainer-stack-utils/client"
	"github.com/greenled/portainer-stack-utils/common"
	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stackAccessCmd represents the stack access command
var stackAccessCmd = &cobra.Command{
	Use:   "access <stackName>",
	Short: "Set access control for stack",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portainerClient, clientRetrievalErr := common.GetClient()
		common.CheckError(clientRetrievalErr)

		stackName := args[0]

		setAdmins := viper.GetBool("stack.access.admins")
		setPrivate := viper.GetBool("stack.access.private")
		setPublic := viper.GetBool("stack.access.public")

		if (setAdmins && setPrivate) || (setPrivate && setPublic) || (setPublic && setAdmins) {
			logrus.Fatal("only one of --admins, --private or --public flags can be used")
		}

		var endpoint portainer.Endpoint
		if endpointName := viper.GetString("stack.access.endpoint"); endpointName == "" {
			// Guess endpoint if not set
			logrus.WithFields(logrus.Fields{
				"implications": "Command will fail if there is not exactly one endpoint available",
			}).Warning("Endpoint not set")
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetDefaultEndpoint()
			common.CheckError(endpointRetrievalErr)
			endpointName = endpoint.Name
			logrus.WithFields(logrus.Fields{
				"endpoint": endpointName,
			}).Debug("Using the only available endpoint")
		} else {
			// Get endpoint by name
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetEndpointByName(endpointName)
			common.CheckError(endpointRetrievalErr)
		}

		logrus.WithFields(logrus.Fields{
			"endpoint": endpoint.Name,
		}).Debug("Getting stack access control info")

		if setAdmins {
			// We are removing an access control
			resourceControl, err := common.GetStackPortainerAccessControl(endpoint.ID, stackName)
			if err == nil {
				err = portainerClient.ResourceControlDelete(resourceControl.ID)
			} else if err != common.ErrAccessControlNotFound {
				common.CheckError(err)
			}
		} else {
			// We may be creating a new access control
			resourceControlCreateOptions := client.ResourceControlCreateOptions{
				ResourceID: stackName,
				Type:       client.ResourceStack,
			}
			if setPrivate {
				currentUser, err := common.GetUserByName(portainerClient.GetUsername())
				common.CheckError(err)
				resourceControlCreateOptions.Users = []portainer.UserID{currentUser.ID}
			}
			if setPublic {
				resourceControlCreateOptions.Public = true
			}
			_, err := portainerClient.ResourceControlCreate(resourceControlCreateOptions)

			if err != nil {
				genericError, isGenericError := err.(*client.GenericError)
				if isGenericError && genericError.Code == http.StatusConflict {
					// We are updating an existing access control
					resourceControl, err := common.GetStackPortainerAccessControl(endpoint.ID, stackName)
					if err == nil {
						resourceControlUpdateOptions := client.ResourceControlUpdateOptions{
							ID: resourceControl.ID,
						}
						if setPrivate {
							currentUser, err := common.GetUserByName(portainerClient.GetUsername())
							common.CheckError(err)
							resourceControlUpdateOptions.Users = []portainer.UserID{currentUser.ID}
						}
						if setPublic {
							resourceControlUpdateOptions.Public = true
						}
						_, err := portainerClient.ResourceControlUpdate(resourceControlUpdateOptions)
						common.CheckError(err)
					} else {
						// Something else happened
						common.CheckError(err)
					}
				} else {
					// Something else happened
					common.CheckError(err)
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"stack": stackName,
		}).Info("Access control set")
	},
}

func init() {
	stackCmd.AddCommand(stackAccessCmd)

	stackAccessCmd.Flags().String("endpoint", "", "Endpoint name.")
	stackAccessCmd.Flags().Bool("admins", false, "Permit access to this stack to administrators only.")
	stackAccessCmd.Flags().Bool("private", false, "Permit access to this stack to the current user only.")
	stackAccessCmd.Flags().Bool("public", false, "Permit access to this stack to any user.")
	viper.BindPFlag("stack.access.endpoint", stackAccessCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.access.admins", stackAccessCmd.Flags().Lookup("admins"))
	viper.BindPFlag("stack.access.private", stackAccessCmd.Flags().Lookup("private"))
	viper.BindPFlag("stack.access.public", stackAccessCmd.Flags().Lookup("public"))
}

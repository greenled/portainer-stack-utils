package common

import (
	"fmt"
	"net/http"

	"github.com/greenled/portainer-stack-utils/client"
	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewAccessCmd creates a new Cobra command for Docker resource access control management
func NewAccessCmd(resourceType client.ResourceType, argumentName string) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("access <%s>", argumentName),
		Short: fmt.Sprintf("Set access control for %s", resourceType),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			portainerClient, clientRetrievalErr := GetClient()
			CheckError(clientRetrievalErr)

			resourceID := args[0]

			setAdmins := viper.GetBool(fmt.Sprintf("%s.access.admins", resourceType))
			setPrivate := viper.GetBool(fmt.Sprintf("%s.access.private", resourceType))
			setPublic := viper.GetBool(fmt.Sprintf("%s.access.public", resourceType))

			if (setAdmins && setPrivate) || (setPrivate && setPublic) || (setPublic && setAdmins) {
				logrus.Fatal("only one of --admins, --private or --public flags can be used")
			}

			var endpoint portainer.Endpoint
			if endpointName := viper.GetString(fmt.Sprintf("%s.access.endpoint", resourceType)); endpointName == "" {
				// Guess endpoint if not set
				logrus.WithFields(logrus.Fields{
					"implications": "Command will fail if there is not exactly one endpoint available",
				}).Warning("Endpoint not set")
				var endpointRetrievalErr error
				endpoint, endpointRetrievalErr = GetDefaultEndpoint()
				CheckError(endpointRetrievalErr)
				endpointName = endpoint.Name
				logrus.WithFields(logrus.Fields{
					"endpoint": endpointName,
				}).Debug("Using the only available endpoint")
			} else {
				// Get endpoint by name
				var endpointRetrievalErr error
				endpoint, endpointRetrievalErr = GetEndpointByName(endpointName)
				CheckError(endpointRetrievalErr)
			}

			logrus.WithFields(logrus.Fields{
				"endpoint": endpoint.Name,
			}).Debug(fmt.Sprintf("Getting %s access control info", resourceType))

			if setAdmins {
				// We are removing an access control
				resourceControl, err := GetDockerResourcePortainerAccessControl(endpoint.ID, resourceID, resourceType)
				if err == nil {
					err = portainerClient.ResourceControlDelete(resourceControl.ID)
				} else if err != ErrAccessControlNotFound {
					CheckError(err)
				}
			} else {
				// We may be creating a new access control
				resourceControlCreateOptions := client.ResourceControlCreateOptions{
					ResourceID: resourceID,
					Type:       resourceType,
				}
				if setPrivate {
					currentUser, err := GetUserByName(portainerClient.GetUsername())
					CheckError(err)
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
						resourceControl, err := GetDockerResourcePortainerAccessControl(endpoint.ID, resourceID, resourceType)
						if err == nil {
							resourceControlUpdateOptions := client.ResourceControlUpdateOptions{
								ID: resourceControl.ID,
							}
							if setPrivate {
								currentUser, err := GetUserByName(portainerClient.GetUsername())
								CheckError(err)
								resourceControlUpdateOptions.Users = []portainer.UserID{currentUser.ID}
							}
							if setPublic {
								resourceControlUpdateOptions.Public = true
							}
							_, err := portainerClient.ResourceControlUpdate(resourceControlUpdateOptions)
							CheckError(err)
						} else {
							// Something else happened
							CheckError(err)
						}
					} else {
						// Something else happened
						CheckError(err)
					}
				}
			}

			logrus.WithFields(logrus.Fields{
				string(resourceType): resourceID,
			}).Info("Access control set")
		},
	}
}

// AccessCmdInitFunc creates an access command for a given Docker resource type
func AccessCmdInitFunc(parentCmd *cobra.Command, resourceControlType client.ResourceType) {
	accessCmd := NewAccessCmd(resourceControlType, fmt.Sprintf("|%sId", resourceControlType))
	parentCmd.AddCommand(accessCmd)

	accessCmd.Flags().String("endpoint", "", "Endpoint name.")
	accessCmd.Flags().Bool("admins", false, fmt.Sprintf("Permit access to this %s to administrators only.", resourceControlType))
	accessCmd.Flags().Bool("private", false, fmt.Sprintf("Permit access to this %s to the current user only.", resourceControlType))
	accessCmd.Flags().Bool("public", false, fmt.Sprintf("Permit access to this %s to any user.", resourceControlType))
	viper.BindPFlag(fmt.Sprintf("%s.access.endpoint", resourceControlType), accessCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag(fmt.Sprintf("%s.access.admins", resourceControlType), accessCmd.Flags().Lookup("admins"))
	viper.BindPFlag(fmt.Sprintf("%s.access.private", resourceControlType), accessCmd.Flags().Lookup("private"))
	viper.BindPFlag(fmt.Sprintf("%s.access.public", resourceControlType), accessCmd.Flags().Lookup("public"))
}

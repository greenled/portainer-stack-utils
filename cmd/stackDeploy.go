package cmd

import (
	"io/ioutil"

	portainer "github.com/portainer/portainer/api"

	"github.com/sirupsen/logrus"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stackDeployCmd represents the undeploy command
var stackDeployCmd = &cobra.Command{
	Use:     "deploy STACK_NAME",
	Short:   "Deploy a new stack or update an existing one",
	Aliases: []string{"up", "create"},
	Example: "  psu stack deploy mystack --stack-file mystack.yml",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var loadedEnvironmentVariables []portainer.Pair
		if viper.GetString("stack.deploy.env-file") != "" {
			var loadingErr error
			loadedEnvironmentVariables, loadingErr = loadEnvironmentVariablesFile(viper.GetString("stack.deploy.env-file"))
			common.CheckError(loadingErr)
		}

		portainerClient, clientRetrievalErr := common.GetClient()
		common.CheckError(clientRetrievalErr)

		stackName := args[0]

		var endpoint portainer.Endpoint
		if endpointName := viper.GetString("stack.deploy.endpoint"); endpointName == "" {
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

		endpointSwarmClusterId, selectionErr := common.GetEndpointSwarmClusterId(endpoint.ID)
		if selectionErr == nil {
			// It's a swarm cluster
		} else if selectionErr == common.ErrStackClusterNotFound {
			// It's not a swarm cluster
		} else {
			// Something else happened
			common.CheckError(selectionErr)
		}

		logrus.WithFields(logrus.Fields{
			"stack":    stackName,
			"endpoint": endpoint.Name,
		}).Debug("Getting stack")
		retrievedStack, stackRetrievalErr := common.GetStackByName(stackName, endpointSwarmClusterId, endpoint.ID)
		if stackRetrievalErr == nil {
			// We are updating an existing stack
			logrus.WithFields(logrus.Fields{
				"stack": retrievedStack.Name,
			}).Debug("Stack found")

			var stackFileContent string
			if viper.GetString("stack.deploy.stack-file") != "" {
				var loadingErr error
				stackFileContent, loadingErr = loadStackFile(viper.GetString("stack.deploy.stack-file"))
				common.CheckError(loadingErr)
			} else {
				var stackFileContentRetrievalErr error
				logrus.WithFields(logrus.Fields{
					"stack": retrievedStack.Name,
				}).Debug("Getting stack file content")
				stackFileContent, stackFileContentRetrievalErr = portainerClient.GetStackFileContent(retrievedStack.ID)
				common.CheckError(stackFileContentRetrievalErr)
			}

			var newEnvironmentVariables []portainer.Pair
			if viper.GetBool("stack.deploy.replace-env") {
				newEnvironmentVariables = loadedEnvironmentVariables
			} else {
				// Merge stack environment variables with the loaded ones
				newEnvironmentVariables = retrievedStack.Env
			LoadedVariablesLoop:
				for _, loadedEnvironmentVariable := range loadedEnvironmentVariables {
					for _, newEnvironmentVariable := range newEnvironmentVariables {
						if loadedEnvironmentVariable.Name == newEnvironmentVariable.Name {
							newEnvironmentVariable.Value = loadedEnvironmentVariable.Value
							continue LoadedVariablesLoop
						}
					}
					newEnvironmentVariables = append(newEnvironmentVariables, portainer.Pair{
						Name:  loadedEnvironmentVariable.Name,
						Value: loadedEnvironmentVariable.Value,
					})
				}
			}

			logrus.WithFields(logrus.Fields{
				"stack": retrievedStack.Name,
			}).Info("Updating stack")
			err := portainerClient.UpdateStack(retrievedStack, newEnvironmentVariables, stackFileContent, viper.GetBool("stack.deploy.prune"), endpoint.ID)
			common.CheckError(err)
		} else if stackRetrievalErr == common.ErrStackNotFound {
			// We are deploying a new stack
			logrus.WithFields(logrus.Fields{
				"stack": stackName,
			}).Debug("Stack not found")

			if viper.GetString("stack.deploy.stack-file") == "" {
				logrus.Fatal(`required flag(s) "stack-file" not set`)
			}
			stackFileContent, loadingErr := loadStackFile(viper.GetString("stack.deploy.stack-file"))
			common.CheckError(loadingErr)

			if endpointSwarmClusterId != "" {
				// It's a swarm cluster
				logrus.WithFields(logrus.Fields{
					"stack":    stackName,
					"endpoint": endpoint.Name,
				}).Info("Creating stack")
				stack, deploymentErr := portainerClient.CreateSwarmStack(stackName, loadedEnvironmentVariables, stackFileContent, endpointSwarmClusterId, endpoint.ID)
				common.CheckError(deploymentErr)
				logrus.WithFields(logrus.Fields{
					"stack":    stack.Name,
					"endpoint": endpoint.Name,
					"id":       stack.ID,
				}).Info("Stack created")
			} else {
				// It's not a swarm cluster
				logrus.WithFields(logrus.Fields{
					"stack":    stackName,
					"endpoint": endpoint.Name,
				}).Info("Creating stack")
				stack, deploymentErr := portainerClient.CreateComposeStack(stackName, loadedEnvironmentVariables, stackFileContent, endpoint.ID)
				common.CheckError(deploymentErr)
				logrus.WithFields(logrus.Fields{
					"stack":    stack.Name,
					"endpoint": endpoint.Name,
					"id":       stack.ID,
				}).Info("Stack created")
			}
		} else {
			// Something else happened
			common.CheckError(stackRetrievalErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackDeployCmd)

	stackDeployCmd.Flags().StringP("stack-file", "c", "", "Path to a file with the content of the stack.")
	stackDeployCmd.Flags().String("endpoint", "", "Endpoint name.")
	stackDeployCmd.Flags().StringP("env-file", "e", "", "Path to a file with environment variables used during stack deployment.")
	stackDeployCmd.Flags().Bool("replace-env", false, "Replace environment variables instead of merging them.")
	stackDeployCmd.Flags().BoolP("prune", "r", false, "Prune services that are no longer referenced (only available for Swarm stacks).")
	viper.BindPFlag("stack.deploy.stack-file", stackDeployCmd.Flags().Lookup("stack-file"))
	viper.BindPFlag("stack.deploy.endpoint", stackDeployCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.deploy.env-file", stackDeployCmd.Flags().Lookup("env-file"))
	viper.BindPFlag("stack.deploy.replace-env", stackDeployCmd.Flags().Lookup("replace-env"))
	viper.BindPFlag("stack.deploy.prune", stackDeployCmd.Flags().Lookup("prune"))
}

func loadStackFile(path string) (string, error) {
	loadedStackFileContentBytes, readingErr := ioutil.ReadFile(path)
	if readingErr != nil {
		return "", readingErr
	}
	return string(loadedStackFileContentBytes), nil
}

// Load environment variables
func loadEnvironmentVariablesFile(path string) ([]portainer.Pair, error) {
	var variables []portainer.Pair
	variablesMap, readingErr := godotenv.Read(path)
	if readingErr != nil {
		return []portainer.Pair{}, readingErr
	}

	for key, value := range variablesMap {
		variables = append(variables, portainer.Pair{
			Name:  key,
			Value: value,
		})
	}

	return variables, nil
}

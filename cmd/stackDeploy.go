package cmd

import (
	"io/ioutil"

	"github.com/greenled/portainer-stack-utils/client"
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
	Example: "psu stack deploy mystack --stack-file mystack.yml",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var loadedEnvironmentVariables []client.StackEnv
		if viper.GetString("stack.deploy.env-file") != "" {
			var loadingErr error
			loadedEnvironmentVariables, loadingErr = loadEnvironmentVariablesFile(viper.GetString("stack.deploy.env-file"))
			common.CheckError(loadingErr)
		}

		portainerClient, clientRetrievalErr := common.GetClient()
		common.CheckError(clientRetrievalErr)

		stackName := args[0]
		logrus.WithFields(logrus.Fields{
			"stack": stackName,
		}).Debug("Getting stack")
		retrievedStack, stackRetrievalErr := common.GetStackByName(stackName)
		switch stackRetrievalErr.(type) {
		case nil:
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
				stackFileContent, stackFileContentRetrievalErr = portainerClient.GetStackFileContent(retrievedStack.Id)
				common.CheckError(stackFileContentRetrievalErr)
			}

			var newEnvironmentVariables []client.StackEnv
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
					newEnvironmentVariables = append(newEnvironmentVariables, client.StackEnv{
						Name:  loadedEnvironmentVariable.Name,
						Value: loadedEnvironmentVariable.Value,
					})
				}
			}

			logrus.WithFields(logrus.Fields{
				"stack": retrievedStack.Name,
			}).Info("Updating stack")
			err := portainerClient.UpdateStack(retrievedStack, newEnvironmentVariables, stackFileContent, viper.GetBool("stack.deploy.prune"), viper.GetString("stack.deploy.endpoint"))
			common.CheckError(err)
		case *common.StackNotFoundError:
			// We are deploying a new stack
			logrus.WithFields(logrus.Fields{
				"stack": stackName,
			}).Debug("Stack not found")

			if viper.GetString("stack.deploy.stack-file") == "" {
				logrus.WithFields(logrus.Fields{
					"flag": "--stack-file",
				}).Fatal("Provide required flag")
			}
			stackFileContent, loadingErr := loadStackFile(viper.GetString("stack.deploy.stack-file"))
			common.CheckError(loadingErr)

			swarmClusterId, selectionErr := getSwarmClusterId()
			endpointId := viper.GetString("stack.deploy.endpoint")
			switch selectionErr.(type) {
			case nil:
				// It's a swarm cluster
				logrus.WithFields(logrus.Fields{
					"stack":    stackName,
					"endpoint": endpointId,
					"cluster":  swarmClusterId,
				}).Info("Creating stack")
				deploymentErr := portainerClient.CreateSwarmStack(stackName, loadedEnvironmentVariables, stackFileContent, swarmClusterId, endpointId)
				common.CheckError(deploymentErr)
			case *valueNotFoundError:
				// It's not a swarm cluster
				logrus.WithFields(logrus.Fields{
					"stack":    stackName,
					"endpoint": endpointId,
				}).Info("Creating stack")
				deploymentErr := portainerClient.CreateComposeStack(stackName, loadedEnvironmentVariables, stackFileContent, endpointId)
				common.CheckError(deploymentErr)
			default:
				// Something else happened
				common.CheckError(stackRetrievalErr)
			}
		default:
			// Something else happened
			common.CheckError(stackRetrievalErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackDeployCmd)

	stackDeployCmd.Flags().StringP("stack-file", "c", "", "Path to a file with the content of the stack.")
	stackDeployCmd.Flags().String("endpoint", "1", "Endpoint ID.")
	stackDeployCmd.Flags().StringP("env-file", "e", "", "Path to a file with environment variables used during stack deployment.")
	stackDeployCmd.Flags().Bool("replace-env", false, "Replace environment variables instead of merging them.")
	stackDeployCmd.Flags().BoolP("prune", "r", false, "Prune services that are no longer referenced (only available for Swarm stacks).")
	viper.BindPFlag("stack.deploy.stack-file", stackDeployCmd.Flags().Lookup("stack-file"))
	viper.BindPFlag("stack.deploy.endpoint", stackDeployCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.deploy.env-file", stackDeployCmd.Flags().Lookup("env-file"))
	viper.BindPFlag("stack.deploy.replace-env", stackDeployCmd.Flags().Lookup("replace-env"))
	viper.BindPFlag("stack.deploy.prune", stackDeployCmd.Flags().Lookup("prune"))
}

func getSwarmClusterId() (id string, err error) {
	// Get docker information for endpoint
	client, err := common.GetClient()
	if err != nil {
		return
	}

	endpointId := viper.GetString("stack.deploy.endpoint")
	logrus.WithFields(logrus.Fields{
		"endpoint": endpointId,
	}).Debug("Getting endpoint's Docker info")
	result, err := client.GetEndpointDockerInfo(endpointId)
	if err != nil {
		return
	}

	// Get swarm (if any) information for endpoint
	swarmClusterId, err := selectValue(result, []string{"Swarm", "Cluster", "ID"})
	if err != nil {
		return
	}
	id = swarmClusterId.(string)

	return
}

func selectValue(jsonMap map[string]interface{}, jsonPath []string) (interface{}, error) {
	value := jsonMap[jsonPath[0]]
	if value == nil {
		return nil, &valueNotFoundError{}
	} else if len(jsonPath) > 1 {
		return selectValue(value.(map[string]interface{}), jsonPath[1:])
	} else {
		return value, nil
	}
}

func loadStackFile(path string) (string, error) {
	loadedStackFileContentBytes, readingErr := ioutil.ReadFile(path)
	if readingErr != nil {
		return "", readingErr
	}
	return string(loadedStackFileContentBytes), nil
}

// Load environment variables
func loadEnvironmentVariablesFile(path string) ([]client.StackEnv, error) {
	var variables []client.StackEnv
	variablesMap, readingErr := godotenv.Read(path)
	if readingErr != nil {
		return []client.StackEnv{}, readingErr
	}

	for key, value := range variablesMap {
		variables = append(variables, client.StackEnv{
			Name:  key,
			Value: value,
		})
	}

	return variables, nil
}

type valueNotFoundError struct{}

func (e *valueNotFoundError) Error() string {
	return "Value not found"
}

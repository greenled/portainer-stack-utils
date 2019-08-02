package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

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
		var loadedEnvironmentVariables []common.StackEnv
		if viper.GetString("stack.deploy.env-file") != "" {
			var loadingErr error
			loadedEnvironmentVariables, loadingErr = loadEnvironmentVariablesFile(viper.GetString("stack.deploy.env-file"))
			common.CheckError(loadingErr)
		}

		client, clientRetrievalErr := common.GetClient()
		common.CheckError(clientRetrievalErr)

		stackName := args[0]
		retrievedStack, stackRetrievalErr := common.GetStackByName(stackName)
		switch stackRetrievalErr.(type) {
		case nil:
			// We are updating an existing stack
			common.PrintVerbose(fmt.Sprintf("Stack %s found. Updating...", retrievedStack.Name))

			var stackFileContent string
			if viper.GetString("stack.deploy.stack-file") != "" {
				var loadingErr error
				stackFileContent, loadingErr = loadStackFile(viper.GetString("stack.deploy.stack-file"))
				common.CheckError(loadingErr)
			} else {
				var stackFileContentRetrievalErr error
				stackFileContent, stackFileContentRetrievalErr = client.GetStackFileContent(retrievedStack.Id)
				common.CheckError(stackFileContentRetrievalErr)
			}

			var newEnvironmentVariables []common.StackEnv
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
					newEnvironmentVariables = append(newEnvironmentVariables, common.StackEnv{
						Name:  loadedEnvironmentVariable.Name,
						Value: loadedEnvironmentVariable.Value,
					})
				}
			}

			err := client.UpdateStack(retrievedStack, newEnvironmentVariables, stackFileContent, viper.GetBool("stack.deploy.prune"), viper.GetString("stack.deploy.endpoint"))
			common.CheckError(err)
		case *common.StackNotFoundError:
			// We are deploying a new stack
			common.PrintVerbose(fmt.Sprintf("Stack %s not found. Deploying...", stackName))

			if viper.GetString("stack.deploy.stack-file") == "" {
				log.Fatalln("Specify a docker-compose file with --stack-file")
			}
			stackFileContent, loadingErr := loadStackFile(viper.GetString("stack.deploy.stack-file"))
			common.CheckError(loadingErr)

			swarmClusterId, selectionErr := getSwarmClusterId()
			switch selectionErr.(type) {
			case nil:
				// It's a swarm cluster
				common.PrintVerbose(fmt.Sprintf("Swarm cluster found with id %s", swarmClusterId))
				deploymentErr := client.CreateSwarmStack(stackName, loadedEnvironmentVariables, stackFileContent, swarmClusterId, viper.GetString("stack.deploy.endpoint"))
				common.CheckError(deploymentErr)
			case *valueNotFoundError:
				// It's not a swarm cluster
				common.PrintVerbose("Swarm cluster not found")
				deploymentErr := client.CreateComposeStack(stackName, loadedEnvironmentVariables, stackFileContent, viper.GetString("stack.deploy.endpoint"))
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

	stackDeployCmd.Flags().StringP("stack-file", "c", "", "path to a file with the content of the stack")
	stackDeployCmd.Flags().String("endpoint", "1", "endpoint ID")
	stackDeployCmd.Flags().StringP("env-file", "e", "", "path to a file with environment variables used during stack deployment")
	stackDeployCmd.Flags().Bool("replace-env", false, "replace environment variables instead of merging them")
	stackDeployCmd.Flags().BoolP("prune", "r", false, "prune services that are no longer referenced (only available for Swarm stacks)")
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

	result, err := client.GetEndpointDockerInfo(viper.GetString("stack.deploy.endpoint"))
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
func loadEnvironmentVariablesFile(path string) ([]common.StackEnv, error) {
	var variables []common.StackEnv
	variablesMap, readingErr := godotenv.Read(path)
	if readingErr != nil {
		return []common.StackEnv{}, readingErr
	}

	for key, value := range variablesMap {
		variables = append(variables, common.StackEnv{
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

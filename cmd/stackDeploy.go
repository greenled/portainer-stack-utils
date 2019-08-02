package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
				stackFileContent, stackFileContentRetrievalErr = getStackFileContent(retrievedStack.Id)
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

			updateErr := updateStack(retrievedStack, newEnvironmentVariables, stackFileContent, viper.GetBool("stack.deploy.prune"))
			common.CheckError(updateErr)
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
				deploymentErr := deploySwarmStack(stackName, loadedEnvironmentVariables, stackFileContent, swarmClusterId)
				common.CheckError(deploymentErr)
			case *valueNotFoundError:
				// It's not a swarm cluster
				common.PrintVerbose("Swarm cluster not found")
				deploymentErr := deployComposeStack(stackName, loadedEnvironmentVariables, stackFileContent)
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

func deploySwarmStack(stackName string, environmentVariables []common.StackEnv, dockerComposeFileContent string, swarmClusterId string) (err error) {
	client, err := common.GetClient()
	if err != nil {
		return
	}

	reqBody := common.StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		SwarmID:          swarmClusterId,
		StackFileContent: dockerComposeFileContent,
	}

	err = client.DoJSON(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%s", 1, "string", viper.GetString("stack.deploy.endpoint")), http.MethodPost, &reqBody, nil)

	return
}

func deployComposeStack(stackName string, environmentVariables []common.StackEnv, stackFileContent string) (err error) {
	client, err := common.GetClient()
	if err != nil {
		return
	}

	reqBody := common.StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
	}

	err = client.DoJSON(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%s", 2, "string", viper.GetString("stack.deploy.endpoint")), http.MethodPost, &reqBody, nil)

	return
}

func updateStack(stack common.Stack, environmentVariables []common.StackEnv, stackFileContent string, prune bool) (err error) {
	client, err := common.GetClient()
	if err != nil {
		return
	}

	reqBody := common.StackUpdateRequest{
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
		Prune:            prune,
	}

	err = client.DoJSON(fmt.Sprintf("stacks/%v?endpointId=%s", stack.Id, viper.GetString("stack.deploy.endpoint")), http.MethodPut, &reqBody, nil)

	return
}

func getSwarmClusterId() (id string, err error) {
	// Get docker information for endpoint
	client, err := common.GetClient()
	if err != nil {
		return
	}

	var result map[string]interface{}

	err = client.DoJSON(fmt.Sprintf("endpoints/%v/docker/info", viper.GetString("stack.deploy.endpoint")), http.MethodGet, nil, &result)
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

func getStackFileContent(stackId uint32) (content string, err error) {
	client, err := common.GetClient()
	if err != nil {
		return
	}

	var respBody common.StackFileInspectResponse

	err = client.DoJSON(fmt.Sprintf("stacks/%v/file", stackId), http.MethodGet, nil, respBody)
	if err != nil {
		return
	}

	content = respBody.StackFileContent

	return
}

type valueNotFoundError struct{}

func (e *valueNotFoundError) Error() string {
	return "Value not found"
}

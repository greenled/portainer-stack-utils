package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/greenled/portainer-stack-utils/common"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
				log.Fatalln(selectionErr)
			}
		default:
			// Something else happened
			log.Fatalln(stackRetrievalErr)
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

func deploySwarmStack(stackName string, environmentVariables []common.StackEnv, dockerComposeFileContent string, swarmClusterId string) error {
	reqBody := common.StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		SwarmID:          swarmClusterId,
		StackFileContent: dockerComposeFileContent,
	}

	reqBodyBytes, marshalingErr := json.Marshal(reqBody)
	if marshalingErr != nil {
		return marshalingErr
	}

	reqUrl, parsingErr := url.Parse(fmt.Sprintf("%s/api/stacks?type=%v&method=%s&endpointId=%s", viper.GetString("url"), 1, "string", viper.GetString("stack.deploy.endpoint")))
	if parsingErr != nil {
		return parsingErr
	}

	req, newRequestErr := http.NewRequest(http.MethodPost, reqUrl.String(), bytes.NewBuffer(reqBodyBytes))
	if newRequestErr != nil {
		return newRequestErr
	}
	headerErr := common.AddAuthorizationHeader(req)
	req.Header.Add("Content-Type", "application/json")
	if headerErr != nil {
		return headerErr
	}
	common.PrintDebugRequest("Deploy stack request", req)

	client := common.NewHttpClient()

	resp, requestExecutionErr := client.Do(req)
	if requestExecutionErr != nil {
		return requestExecutionErr
	}
	common.PrintDebugResponse("Deploy stack response", resp)

	responseErr := common.CheckResponseForErrors(resp)
	if responseErr != nil {
		return responseErr
	}

	return nil
}

func deployComposeStack(stackName string, environmentVariables []common.StackEnv, stackFileContent string) error {
	reqBody := common.StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
	}

	reqBodyBytes, marshalingErr := json.Marshal(reqBody)
	if marshalingErr != nil {
		return marshalingErr
	}

	reqUrl, parsingErr := url.Parse(fmt.Sprintf("%s/api/stacks?type=%v&method=%s&endpointId=%s", viper.GetString("url"), 2, "string", viper.GetString("stack.deploy.endpoint")))
	if parsingErr != nil {
		return parsingErr
	}

	req, newRequestErr := http.NewRequest(http.MethodPost, reqUrl.String(), bytes.NewBuffer(reqBodyBytes))
	if newRequestErr != nil {
		return newRequestErr
	}
	headerErr := common.AddAuthorizationHeader(req)
	req.Header.Add("Content-Type", "application/json")
	if headerErr != nil {
		return headerErr
	}
	common.PrintDebugRequest("Deploy stack request", req)

	client := common.NewHttpClient()

	resp, requestExecutionErr := client.Do(req)
	if requestExecutionErr != nil {
		return requestExecutionErr
	}
	common.PrintDebugResponse("Deploy stack response", resp)

	responseErr := common.CheckResponseForErrors(resp)
	if responseErr != nil {
		return responseErr
	}

	return nil
}

func updateStack(stack common.Stack, environmentVariables []common.StackEnv, stackFileContent string, prune bool) error {
	reqBody := common.StackUpdateRequest{
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
		Prune:            prune,
	}

	reqBodyBytes, marshalingErr := json.Marshal(reqBody)
	if marshalingErr != nil {
		return marshalingErr
	}

	reqUrl, parsingErr := url.Parse(fmt.Sprintf("%s/api/stacks/%v?endpointId=%s", viper.GetString("url"), stack.Id, viper.GetString("stack.deploy.endpoint")))
	if parsingErr != nil {
		return parsingErr
	}

	req, newRequestErr := http.NewRequest(http.MethodPut, reqUrl.String(), bytes.NewBuffer(reqBodyBytes))
	if newRequestErr != nil {
		return newRequestErr
	}
	headerErr := common.AddAuthorizationHeader(req)
	req.Header.Add("Content-Type", "application/json")
	if headerErr != nil {
		return headerErr
	}
	common.PrintDebugRequest("Update stack request", req)

	client := common.NewHttpClient()

	resp, requestExecutionErr := client.Do(req)
	if requestExecutionErr != nil {
		return requestExecutionErr
	}
	common.PrintDebugResponse("Update stack response", resp)

	responseErr := common.CheckResponseForErrors(resp)
	if responseErr != nil {
		return responseErr
	}

	return nil
}

func getSwarmClusterId() (string, error) {
	// Get docker information for endpoint
	reqUrl, parsingErr := url.Parse(fmt.Sprintf("%s/api/endpoints/%v/docker/info", viper.GetString("url"), viper.GetString("stack.deploy.endpoint")))
	if parsingErr != nil {
		return "", parsingErr
	}

	req, newRequestErr := http.NewRequest(http.MethodGet, reqUrl.String(), nil)
	if newRequestErr != nil {
		return "", newRequestErr
	}
	headerErr := common.AddAuthorizationHeader(req)
	if headerErr != nil {
		return "", headerErr
	}
	common.PrintDebugRequest("Get docker info request", req)

	client := common.NewHttpClient()

	resp, requestExecutionErr := client.Do(req)
	if requestExecutionErr != nil {
		return "", requestExecutionErr
	}
	common.PrintDebugResponse("Get docker info response", resp)

	responseErr := common.CheckResponseForErrors(resp)
	if responseErr != nil {
		return "", responseErr
	}

	// Get swarm (if any) information for endpoint
	var result map[string]interface{}
	decodingError := json.NewDecoder(resp.Body).Decode(&result)
	if decodingError != nil {
		return "", decodingError
	}

	swarmClusterId, selectionErr := selectValue(result, []string{"Swarm", "Cluster", "ID"})
	if selectionErr != nil {
		return "", selectionErr
	}

	return swarmClusterId.(string), nil
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

func getStackFileContent(stackId uint32) (string, error) {
	reqUrl, parsingErr := url.Parse(fmt.Sprintf("%s/api/stacks/%v/file", viper.GetString("url"), stackId))
	if parsingErr != nil {
		return "", parsingErr
	}

	req, newRequestErr := http.NewRequest(http.MethodGet, reqUrl.String(), nil)
	if newRequestErr != nil {
		return "", newRequestErr
	}
	headerErr := common.AddAuthorizationHeader(req)
	if headerErr != nil {
		return "", headerErr
	}
	common.PrintDebugRequest("Get stack file content request", req)

	client := common.NewHttpClient()

	resp, requestExecutionErr := client.Do(req)
	if requestExecutionErr != nil {
		return "", requestExecutionErr
	}
	common.PrintDebugResponse("Get stack file content response", resp)

	responseErr := common.CheckResponseForErrors(resp)
	if responseErr != nil {
		return "", responseErr
	}

	var respBody common.StackFileInspectResponse
	decodingErr := json.NewDecoder(resp.Body).Decode(&respBody)
	if decodingErr != nil {
		return "", decodingErr
	}

	return respBody.StackFileContent, nil
}

type valueNotFoundError struct{}

func (e *valueNotFoundError) Error() string {
	return "Value not found"
}

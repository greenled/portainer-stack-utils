package common

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

func GetAllStacks() ([]Stack, error) {
	return GetAllStacksFiltered(StackListFilter{})
}

func GetAllStacksFiltered(filter StackListFilter) ([]Stack, error) {
	PrintVerbose("Getting all stacks...")

	filterJsonBytes, _ := json.Marshal(filter)
	filterJsonString := string(filterJsonBytes)

	reqUrl, err := url.Parse(fmt.Sprintf("%s/api/stacks?filters=%s", viper.GetString("url"), filterJsonString))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	headerErr := AddAuthorizationHeader(req)
	if headerErr != nil {
		return nil, err
	}
	PrintDebugRequest("Get stacks request", req)

	client := NewHttpClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	PrintDebugResponse("Get stacks response", resp)

	CheckError(CheckResponseForErrors(resp))

	var respBody []Stack
	decodingErr := json.NewDecoder(resp.Body).Decode(&respBody)
	CheckError(decodingErr)

	return respBody, nil
}

func GetStackByName(name string) (Stack, error) {
	stacks, err := GetAllStacks()
	if err != nil {
		return Stack{}, err
	}

	PrintVerbose(fmt.Sprintf("Getting stack %s...", name))
	for _, stack := range stacks {
		if stack.Name == name {
			return stack, nil
		}
	}
	return Stack{}, &StackNotFoundError{
		StackName: name,
	}
}

type StackListFilter struct {
	SwarmId    string `json:",omitempty"`
	EndpointId uint32 `json:",omitempty"`
}

// Custom customerrors
type StackNotFoundError struct {
	StackName string
}

func (e *StackNotFoundError) Error() string {
	return fmt.Sprintf("Stack %s not found", e.StackName)
}

func GetAllEndpoints() ([]EndpointSubset, error) {
	PrintVerbose("Getting all endpoints...")

	reqUrl, err := url.Parse(fmt.Sprintf("%s/api/endpoints", viper.GetString("url")))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	headerErr := AddAuthorizationHeader(req)
	if headerErr != nil {
		return nil, err
	}
	PrintDebugRequest("Get endpoints request", req)

	client := NewHttpClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	PrintDebugResponse("Get endpoints response", resp)

	CheckError(CheckResponseForErrors(resp))

	var respBody []EndpointSubset
	decodingErr := json.NewDecoder(resp.Body).Decode(&respBody)
	CheckError(decodingErr)

	return respBody, nil
}

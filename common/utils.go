package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAllStacks() ([]Stack, error) {
	return GetAllStacksFiltered(StackListFilter{})
}

func GetAllStacksFiltered(filter StackListFilter) (stacks []Stack, err error) {
	PrintVerbose("Getting all stacks...")

	client, err := GetClient()
	if err != nil {
		return
	}

	filterJsonBytes, _ := json.Marshal(filter)
	filterJsonString := string(filterJsonBytes)

	err = client.DoJSON(fmt.Sprintf("stacks?filters=%s", filterJsonString), http.MethodGet, nil, &stacks)
	return
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

func GetAllEndpoints() (endpoints []EndpointSubset, err error) {
	PrintVerbose("Getting all endpoints...")

	client, err := GetClient()
	if err != nil {
		return
	}

	err = client.DoJSON("endpoints", http.MethodGet, nil, &endpoints)
	return
}

package common

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/greenled/portainer-stack-utils/client"
)

func GetStackByName(name string, swarmId string, endpointId uint32) (stack client.Stack, err error) {
	portainerClient, err := GetClient()
	if err != nil {
		return
	}

	stacks, err := portainerClient.GetStacks(swarmId, endpointId)
	if err != nil {
		return
	}

	for _, stack := range stacks {
		if stack.Name == name {
			return stack, nil
		}
	}
	err = &StackNotFoundError{
		StackName: name,
	}
	return
}

func GetEndpointSwarmClusterId(endpointId uint32) (endpointSwarmClusterId string, err error) {
	// Get docker information for endpoint
	portainerClient, err := GetClient()
	if err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{
		"endpoint": endpointId,
	}).Debug("Getting endpoint's Docker info")
	result, err := portainerClient.GetEndpointDockerInfo(endpointId)
	if err != nil {
		return
	}

	// Get swarm (if any) information for endpoint
	id, selectionErr := selectValue(result, []string{"Swarm", "Cluster", "ID"})
	switch selectionErr.(type) {
	case nil:
		endpointSwarmClusterId = id.(string)
	case *valueNotFoundError:
		err = &StackClusterNotFoundError{}
	default:
		err = selectionErr
	}

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

// Custom customerrors
type StackNotFoundError struct {
	StackName string
}

func (e *StackNotFoundError) Error() string {
	return fmt.Sprintf("Stack %s not found", e.StackName)
}

type valueNotFoundError struct{}

func (e *valueNotFoundError) Error() string {
	return "Value not found"
}

type StackClusterNotFoundError struct{}

func (e *StackClusterNotFoundError) Error() string {
	return "Stack cluster not found"
}

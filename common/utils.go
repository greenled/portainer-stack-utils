package common

import (
	"errors"
	"fmt"
	"reflect"

	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
)

func GetDefaultEndpoint() (endpoint portainer.Endpoint, err error) {
	portainerClient, err := GetClient()
	if err != nil {
		return
	}

	logrus.Debug("Getting endpoints")
	endpoints, err := portainerClient.GetEndpoints()
	if err != nil {
		return
	}

	if len(endpoints) == 0 {
		err = errors.New("No endpoints available")
		return
	} else if len(endpoints) > 1 {
		err = errors.New("Several endpoints available")
		return
	}
	endpoint = endpoints[0]

	return
}

func GetStackByName(name string, swarmId string, endpointId portainer.EndpointID) (stack portainer.Stack, err error) {
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

func GetEndpointSwarmClusterId(endpointId portainer.EndpointID) (endpointSwarmClusterId string, err error) {
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

func GetFormatHelp(v interface{}) (r string) {
	typeOfV := reflect.TypeOf(v)
	r = fmt.Sprintf(`
Format:
  The --format flag accepts a Go template, which is passed a %s.%s object:

%s
`, typeOfV.PkgPath(), typeOfV.Name(), fmt.Sprintf("%s%s", "  ", repr(typeOfV, "  ", "  ")))
	return
}

func repr(t reflect.Type, margin, beforeMargin string) (r string) {
	switch t.Kind() {
	case reflect.Struct:
		r = fmt.Sprintln("{")
		for i := 0; i < t.NumField(); i++ {
			tField := t.Field(i)
			r += fmt.Sprintln(fmt.Sprintf("%s%s%s %s", beforeMargin, margin, tField.Name, repr(tField.Type, margin, beforeMargin+margin)))
		}
		r += fmt.Sprintf("%s}", beforeMargin)
	case reflect.Array, reflect.Slice:
		r = fmt.Sprintf("[]%s", repr(t.Elem(), margin, beforeMargin))
	default:
		r = fmt.Sprintf("%s", t.Name())
	}
	return
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

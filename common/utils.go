package common

import (
	"fmt"
	"reflect"

	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
)

const (
	ErrStackNotFound             = Error("Stack not found")
	ErrStackClusterNotFound      = Error("Stack cluster not found")
	ErrEndpointNotFound          = Error("Endpoint not found")
	ErrEndpointGroupNotFound     = Error("Endpoint group not found")
	ErrSeveralEndpointsAvailable = Error("Several endpoints available")
	ErrNoEndpointsAvailable      = Error("No endpoints available")
)

const (
	valueNotFoundError = Error("Value not found")
)

// Error represents an application error.
type Error string

// Error returns the error message.
func (e Error) Error() string {
	return string(e)
}

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
		err = ErrNoEndpointsAvailable
		return
	} else if len(endpoints) > 1 {
		err = ErrSeveralEndpointsAvailable
		return
	}
	endpoint = endpoints[0]

	return
}

func GetStackByName(name string, swarmID string, endpointID portainer.EndpointID) (stack portainer.Stack, err error) {
	portainerClient, err := GetClient()
	if err != nil {
		return
	}

	stacks, err := portainerClient.GetStacks(swarmID, endpointID)
	if err != nil {
		return
	}

	for _, stack := range stacks {
		if stack.Name == name {
			return stack, nil
		}
	}
	err = ErrStackNotFound
	return
}

func GetEndpointByName(name string) (endpoint portainer.Endpoint, err error) {
	portainerClient, err := GetClient()
	if err != nil {
		return
	}

	endpoints, err := portainerClient.GetEndpoints()
	if err != nil {
		return
	}

	for _, endpoint := range endpoints {
		if endpoint.Name == name {
			return endpoint, nil
		}
	}
	err = ErrEndpointNotFound
	return
}

func GetEndpointGroupByName(name string) (endpointGroup portainer.EndpointGroup, err error) {
	portainerClient, err := GetClient()
	if err != nil {
		return
	}

	endpointGroups, err := portainerClient.GetEndpointGroups()
	if err != nil {
		return
	}

	for _, endpointGroup := range endpointGroups {
		if endpointGroup.Name == name {
			return endpointGroup, nil
		}
	}
	err = ErrEndpointGroupNotFound
	return
}

func GetEndpointFromListByID(endpoints []portainer.Endpoint, id portainer.EndpointID) (endpoint portainer.Endpoint, err error) {
	for i := range endpoints {
		if endpoints[i].ID == id {
			return endpoints[i], err
		}
	}
	return endpoint, ErrEndpointNotFound
}

func GetEndpointFromListByName(endpoints []portainer.Endpoint, name string) (endpoint portainer.Endpoint, err error) {
	for i := range endpoints {
		if endpoints[i].Name == name {
			return endpoints[i], err
		}
	}
	return endpoint, ErrEndpointNotFound
}

func GetEndpointSwarmClusterID(endpointID portainer.EndpointID) (endpointSwarmClusterID string, err error) {
	// Get docker information for endpoint
	portainerClient, err := GetClient()
	if err != nil {
		return
	}

	result, err := portainerClient.GetEndpointDockerInfo(endpointID)
	if err != nil {
		return
	}

	// Get swarm (if any) information for endpoint
	id, selectionErr := selectValue(result, []string{"Swarm", "Cluster", "ID"})
	if selectionErr == nil {
		endpointSwarmClusterID = id.(string)
	} else if selectionErr == valueNotFoundError {
		err = ErrStackClusterNotFound
	} else {
		err = selectionErr
	}

	return
}

func selectValue(jsonMap map[string]interface{}, jsonPath []string) (interface{}, error) {
	value := jsonMap[jsonPath[0]]
	if value == nil {
		return nil, valueNotFoundError
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

package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

// StackCreateComposeOptions represents options passed to PortainerClient.StackCreateCompose()
type StackCreateComposeOptions struct {
	StackName            string
	EnvironmentVariables []portainer.Pair
	StackFileContent     string
	EndpointID           portainer.EndpointID
}

// StackCreateRequest represents the body of a request to POST /stacks
type StackCreateRequest struct {
	Name             string
	SwarmID          string
	StackFileContent string
	Env              []portainer.Pair `json:",omitempty"`
}

// StackCreateSwarmOptions represents options passed to PortainerClient.StackCreateSwarm()
type StackCreateSwarmOptions struct {
	StackName            string
	EnvironmentVariables []portainer.Pair
	StackFileContent     string
	SwarmClusterID       string
	EndpointID           portainer.EndpointID
}

func (n *portainerClientImp) StackCreateCompose(options StackCreateComposeOptions) (stack portainer.Stack, err error) {
	reqBody := StackCreateRequest{
		Name:             options.StackName,
		Env:              options.EnvironmentVariables,
		StackFileContent: options.StackFileContent,
	}

	err = n.DoJSONWithToken(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%v", 2, "string", options.EndpointID), http.MethodPost, http.Header{}, &reqBody, &stack)
	return
}

func (n *portainerClientImp) StackCreateSwarm(options StackCreateSwarmOptions) (stack portainer.Stack, err error) {
	reqBody := StackCreateRequest{
		Name:             options.StackName,
		Env:              options.EnvironmentVariables,
		SwarmID:          options.SwarmClusterID,
		StackFileContent: options.StackFileContent,
	}

	err = n.DoJSONWithToken(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%v", 1, "string", options.EndpointID), http.MethodPost, http.Header{}, &reqBody, &stack)
	return
}

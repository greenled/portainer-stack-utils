package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

// StackUpdateOptions represents options passed to PortainerClient.StackUpdate()
type StackUpdateOptions struct {
	Stack                portainer.Stack
	EnvironmentVariables []portainer.Pair
	StackFileContent     string
	Prune                bool
	EndpointID           portainer.EndpointID
}

// StackUpdateRequest represents the body of a request to PUT /stacks/{id}
type StackUpdateRequest struct {
	StackFileContent string
	Env              []portainer.Pair `json:",omitempty"`
	Prune            bool
}

func (n *portainerClientImp) StackUpdate(options StackUpdateOptions) (err error) {
	reqBody := StackUpdateRequest{
		Env:              options.EnvironmentVariables,
		StackFileContent: options.StackFileContent,
		Prune:            options.Prune,
	}

	err = n.doJSONWithToken(fmt.Sprintf("stacks/%v?endpointId=%v", options.Stack.ID, options.EndpointID), http.MethodPut, http.Header{}, &reqBody, nil)
	return
}

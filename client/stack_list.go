package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

// StackListFilter represents a filter for a stack list
type StackListFilter struct {
	SwarmID    string               `json:"SwarmId,omitempty"`
	EndpointID portainer.EndpointID `json:"EndpointId,omitempty"`
}

// StackListOptions represents options passed to PortainerClient.StackList()
type StackListOptions struct {
	Filter StackListFilter
}

func (n *portainerClientImp) StackList(options StackListOptions) (stacks []portainer.Stack, err error) {
	filterJSONBytes, _ := json.Marshal(options.Filter)
	filterJSONString := string(filterJSONBytes)

	err = n.DoJSONWithToken(fmt.Sprintf("stacks?filters=%s", filterJSONString), http.MethodGet, http.Header{}, nil, &stacks)
	return
}

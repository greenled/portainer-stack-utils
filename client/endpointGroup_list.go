package client

import (
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) EndpointGroupList() (endpointGroups []portainer.EndpointGroup, err error) {
	err = n.DoJSONWithToken("endpoint_groups", http.MethodGet, http.Header{}, nil, &endpointGroups)
	return
}

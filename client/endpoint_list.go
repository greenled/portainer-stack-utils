package client

import (
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) EndpointList() (endpoints []portainer.Endpoint, err error) {
	err = n.doJSONWithToken("endpoints", http.MethodGet, http.Header{}, nil, &endpoints)
	return
}

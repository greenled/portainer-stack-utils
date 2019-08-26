package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) EndpointDockerInfo(endpointID portainer.EndpointID) (info map[string]interface{}, err error) {
	err = n.doJSONWithToken(fmt.Sprintf("endpoints/%v/docker/info", endpointID), http.MethodGet, http.Header{}, nil, &info)
	return
}

package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) Proxy(endpointID portainer.EndpointID, req *http.Request) (resp *http.Response, err error) {
	return n.doWithToken(fmt.Sprintf("endpoints/%v/docker%s", endpointID, req.RequestURI), req.Method, req.Body, req.Header)
}

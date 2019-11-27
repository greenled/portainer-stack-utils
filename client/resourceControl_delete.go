package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) ResourceControlDelete(resourceControlID portainer.ResourceControlID) (err error) {
	err = n.DoJSONWithToken(fmt.Sprintf("resource_controls/%d", resourceControlID), http.MethodDelete, http.Header{}, nil, nil)
	return
}

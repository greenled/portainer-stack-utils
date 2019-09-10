package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) StackDelete(stackID portainer.StackID) (err error) {
	err = n.DoJSONWithToken(fmt.Sprintf("stacks/%d", stackID), http.MethodDelete, http.Header{}, nil, nil)
	return
}

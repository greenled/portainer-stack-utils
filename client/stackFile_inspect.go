package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

// StackFileInspectResponse represents the body of a response for a request to GET /stack/{id}/file
type StackFileInspectResponse struct {
	StackFileContent string
}

func (n *portainerClientImp) StackFileInspect(stackID portainer.StackID) (content string, err error) {
	var respBody StackFileInspectResponse

	err = n.doJSONWithToken(fmt.Sprintf("stacks/%v/file", stackID), http.MethodGet, http.Header{}, nil, &respBody)
	if err != nil {
		return
	}

	content = respBody.StackFileContent

	return
}

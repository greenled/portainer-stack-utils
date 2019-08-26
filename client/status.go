package client

import (
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) Status() (status portainer.Status, err error) {
	err = n.doJSONWithToken("status", http.MethodGet, http.Header{}, nil, &status)
	return
}

package client

import (
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

func (n *portainerClientImp) UserList() (users []portainer.User, err error) {
	err = n.doJSONWithToken("users", http.MethodGet, http.Header{}, nil, &users)
	return
}

package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

// ResourceControlUpdateOptions represents options passed to PortainerClient.ResourceControlUpdate()
type ResourceControlUpdateOptions struct {
	ID     portainer.ResourceControlID
	Public bool
	Users  []portainer.UserID
	Teams  []portainer.TeamID
}

// ResourceControlUpdateRequest represents the body of a request to PUT /resource_controls/{id}
type ResourceControlUpdateRequest struct {
	Public bool               `json:",omitempty"`
	Users  []portainer.UserID `json:",omitempty"`
	Teams  []portainer.TeamID `json:",omitempty"`
}

func (n *portainerClientImp) ResourceControlUpdate(options ResourceControlUpdateOptions) (resourceControl portainer.ResourceControl, err error) {
	reqBody := ResourceControlUpdateRequest{
		Public: options.Public,
		Users:  options.Users,
		Teams:  options.Teams,
	}

	err = n.DoJSONWithToken(fmt.Sprintf("resource_controls/%d", options.ID), http.MethodPut, http.Header{}, &reqBody, &resourceControl)
	return
}

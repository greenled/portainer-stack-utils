package client

import (
	"fmt"
	"net/http"

	portainer "github.com/portainer/portainer/api"
)

// ResourceControlCreateOptions represents options passed to PortainerClient.ResourceControlCreate()
type ResourceControlCreateOptions struct {
	ResourceID     string
	Type           ResourceType
	Public         bool
	Users          []portainer.UserID
	Teams          []portainer.TeamID
	SubResourceIDs []string
}

// ResourceControlCreateRequest represents the body of a request to POST /resource_controls
type ResourceControlCreateRequest struct {
	ResourceID     string
	Type           ResourceType
	Public         bool               `json:",omitempty"`
	Users          []portainer.UserID `json:",omitempty"`
	Teams          []portainer.TeamID `json:",omitempty"`
	SubResourceIDs []string           `json:",omitempty"`
}

func (n *portainerClientImp) ResourceControlCreate(options ResourceControlCreateOptions) (resourceControl portainer.ResourceControl, err error) {
	reqBody := ResourceControlCreateRequest{
		ResourceID:     options.ResourceID,
		Type:           options.Type,
		Public:         options.Public,
		Users:          options.Users,
		Teams:          options.Teams,
		SubResourceIDs: options.SubResourceIDs,
	}

	err = n.DoJSONWithToken(fmt.Sprintf("resource_controls"), http.MethodPost, http.Header{}, &reqBody, &resourceControl)
	return
}

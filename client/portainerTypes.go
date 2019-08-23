package client

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"
)

// GetTranslatedStackType returns a stack's Type field (int) translated to it's human readable form (string)
func GetTranslatedStackType(s portainer.Stack) string {
	switch s.Type {
	case 1:
		return "swarm"
	case 2:
		return "compose"
	default:
		return ""
	}
}

// StackCreateRequest represents the body of a request to POST /stacks
type StackCreateRequest struct {
	Name             string
	SwarmID          string
	StackFileContent string
	Env              []portainer.Pair `json:",omitempty"`
}

// StackUpdateRequest represents the body of a request to PUT /stacks/{id}
type StackUpdateRequest struct {
	StackFileContent string
	Env              []portainer.Pair `json:",omitempty"`
	Prune            bool
}

// StackFileInspectResponse represents the body of a response for a request to GET /stack/{id}/file
type StackFileInspectResponse struct {
	StackFileContent string
}

// GenericError represents the body of a generic error returned by the Portainer API
type GenericError struct {
	Err     string
	Details string
}

func (e *GenericError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Err, e.Details)
	} else {
		return fmt.Sprintf("%s", e.Err)
	}
}

// AuthenticateUserRequest represents the body of a request to POST /auth
type AuthenticateUserRequest struct {
	Username string
	Password string
}

// AuthenticateUserResponse represents the body of a response for a request to POST /auth
type AuthenticateUserResponse struct {
	Jwt string
}

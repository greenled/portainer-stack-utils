package client

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"
)

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

type StackCreateRequest struct {
	Name             string
	SwarmID          string
	StackFileContent string
	Env              []portainer.Pair `json:",omitempty"`
}

type StackUpdateRequest struct {
	StackFileContent string
	Env              []portainer.Pair `json:",omitempty"`
	Prune            bool
}

type StackFileInspectResponse struct {
	StackFileContent string
}

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

type AuthenticateUserRequest struct {
	Username string
	Password string
}

type AuthenticateUserResponse struct {
	Jwt string
}

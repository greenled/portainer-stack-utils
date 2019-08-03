package client

import "fmt"

type Stack struct {
	// In the API documentation this field is a String,
	// but it's returned as a number
	Id          uint32
	Name        string
	Type        uint8 // 1 for a Swarm stack, 2 for a Compose stack
	EndpointID  uint
	EntryPoint  string
	SwarmID     string
	ProjectPath string
	Env         []StackEnv
}

func (s *Stack) GetTranslatedStackType() string {
	switch s.Type {
	case 1:
		return "swarm"
	case 2:
		return "compose"
	default:
		return ""
	}
}

type StackEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EndpointSubset struct {
	Id        uint32
	Name      string
	Type      uint8
	URL       string
	PublicURL string
	GroupID   uint32
}

type StackCreateRequest struct {
	Name             string
	SwarmID          string
	StackFileContent string
	Env              []StackEnv `json:",omitempty"`
}

type StackUpdateRequest struct {
	StackFileContent string
	Env              []StackEnv `json:",omitempty"`
	Prune            bool
}

type StackFileInspectResponse struct {
	StackFileContent string
}

type Status struct {
	Authentication     bool
	EndpointManagement bool
	Analytics          bool
	Version            string
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

package client

import (
	portainer "github.com/portainer/portainer/api"
)

// GetTranslatedStackType returns a stack's Type field (int) translated to it's human readable form (string)
func GetTranslatedStackType(t portainer.StackType) string {
	switch t {
	case portainer.DockerSwarmStack:
		return "swarm"
	case portainer.DockerComposeStack:
		return "compose"
	default:
		return ""
	}
}

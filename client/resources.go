package client

type (
	// ResourceType represents a type of Docker or Portainer resource
	ResourceType string
)

const (
	// ResourceContainer represents a Docker resource of type container
	ResourceContainer = ResourceType("container")

	// ResourceService represents a Docker resource of type service
	ResourceService = ResourceType("service")

	// ResourceVolume represents a Docker resource of type volume
	ResourceVolume = ResourceType("volume")

	// ResourceNetwork represents a Docker resource of type network
	ResourceNetwork = ResourceType("network")

	// ResourceSecret represents a Docker resource of type secret
	ResourceSecret = ResourceType("secret")

	// ResourceConfig represents a Docker resource of type config
	ResourceConfig = ResourceType("config")

	// ResourceStack represents a Portainer resource of type stack
	// A Portainer stack is pretty much like a Docker stack, but not the same
	ResourceStack = ResourceType("stack")
)

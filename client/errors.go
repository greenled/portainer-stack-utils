package client

import "fmt"

// GenericError represents the body of a generic error returned by the Portainer API
type GenericError struct {
	Err     string
	Details string
}

func (e *GenericError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Err, e.Details)
	}
	return fmt.Sprintf("%s", e.Err)
}

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

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

// Check if an http.Response object has errors
func checkResponseForErrors(resp *http.Response) error {
	if 300 <= resp.StatusCode {
		// Guess it's a GenericError
		respBody := GenericError{}
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			// It's not a GenericError
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				return err
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return errors.New(string(bodyBytes))
		}
		return &respBody
	}
	return nil
}

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GenericError represents the body of a generic error returned by the Portainer API
type GenericError struct {
	Code    int
	Err     string
	Details string
}

func (e GenericError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Err, e.Details)
	}
	return fmt.Sprintf("%s", e.Err)
}

// Get an http.Response's error (if any)
func getResponseHTTPError(resp *http.Response) error {
	if resp.StatusCode < 300 {
		// There is no error
		return nil
	}

	switch resp.StatusCode {
	// Error codes found in the Portainer API 1.22.0 Swagger spec
	case http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound, http.StatusConflict, http.StatusInternalServerError, http.StatusServiceUnavailable:
		// Guess it's a GenericError
		genericError, err := getResponseGenericHTTPError(resp)
		if err != nil {
			// It's not a GenericError
			return getResponseNonGenericHTTPError(resp)
		}
		return &genericError
	default:
		return getResponseNonGenericHTTPError(resp)
	}
}

func getResponseGenericHTTPError(resp *http.Response) (genericError GenericError, err error) {
	genericError = GenericError{
		Code: resp.StatusCode,
	}
	err = json.NewDecoder(resp.Body).Decode(&genericError)
	return
}

func getResponseNonGenericHTTPError(resp *http.Response) error {
	bodyString, err := getResponseBodyAsString(resp)
	if err != nil {
		return err
	}
	return errors.New(bodyString)
}

func getResponseBodyAsString(resp *http.Response) (bodyString string, err error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	bodyString = string(bodyBytes)
	resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	return
}

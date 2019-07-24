package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// CheckError checks if an error occurred (it's not nil)
func CheckError(err error) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error: %s", err.Error()))
	}
}

func CheckResponseForErrors(resp *http.Response) error {
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

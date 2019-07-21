package common

import (
	"encoding/json"
	"log"
	"net/http"
)

// CheckError checks if an error occurred (it's not nil)
func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckResponseForErrors(resp *http.Response) error {
	if 300 <= resp.StatusCode {
		respBody := GenericError{}
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return err
		}
		return &respBody
	}
	return nil
}

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

var cachedAuthenticationToken string

func GetAuthenticationToken() (string, error) {
	if cachedAuthenticationToken == "" {
		var authenticationTokenRetrievalErr error
		cachedAuthenticationToken, authenticationTokenRetrievalErr = GetNewAuthenticationToken()
		if authenticationTokenRetrievalErr != nil {
			return "", authenticationTokenRetrievalErr
		}
	}
	return cachedAuthenticationToken, nil
}

func GetNewAuthenticationToken() (string, error) {
	PrintVerbose("Getting auth token...")

	reqBody := AuthenticateUserRequest{
		Username: viper.GetString("user"),
		Password: viper.GetString("password"),
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	reqUrl, err := url.Parse(fmt.Sprintf("%s/api/auth", viper.GetString("url")))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl.String(), bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return "", err
	}
	PrintDebugRequest("Get auth token request", req)

	client := NewHttpClient()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	PrintDebugResponse("Get auth token response", resp)

	respErr := CheckResponseForErrors(resp)
	if respErr != nil {
		return "", err
	}

	respBody := AuthenticateUserResponse{}
	decodingErr := json.NewDecoder(resp.Body).Decode(&respBody)
	CheckError(decodingErr)
	PrintDebug(fmt.Sprintf("Auth token: %s", respBody.Jwt))
	return respBody.Jwt, nil
}

func AddAuthorizationHeader(request *http.Request) error {
	token, err := GetAuthenticationToken()
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	return nil
}

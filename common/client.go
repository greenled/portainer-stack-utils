package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var client *PortainerClient

type ClientConfig struct {
	Url           string
	User          string
	Password      string
	Token         string
	DoNotUseToken bool
	Insecure      bool
	Timeout       time.Duration
}

type PortainerClient struct {
	http.Client
	url           *url.URL
	user          string
	password      string
	token         string
	doNotUseToken bool
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

// Do an http request
func (n *PortainerClient) do(uri, method string, request io.Reader, requestType string, headers http.Header) (resp *http.Response, err error) {
	requestUrl, err := n.url.Parse(uri)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, requestUrl.String(), request)
	if err != nil {
		return
	}

	if headers != nil {
		req.Header = headers
	}

	if request != nil {
		req.Header.Set("Content-Type", requestType)
	}

	if !n.doNotUseToken {
		if n.token == "" {
			clientClone, cloneErr := n.Clone()
			if cloneErr != nil {
				return resp, cloneErr
			}
			clientClone.doNotUseToken = true
			n.token, err = clientClone.Authenticate()
			if err != nil {
				return
			}
			PrintDebug(fmt.Sprintf("Auth token: %s", n.token))
		}
		req.Header.Set("Authorization", "Bearer "+n.token)
	}

	PrintDebugRequest("Request", req)

	resp, err = n.Do(req)
	if err != nil {
		return
	}

	err = checkResponseForErrors(resp)
	if err != nil {
		return
	}

	PrintDebugResponse("Response", resp)

	return
}

// Do a JSON http request
func (n *PortainerClient) doJSON(uri, method string, request interface{}, response interface{}) error {
	var body io.Reader

	if request != nil {
		reqBodyBytes, err := json.Marshal(request)
		if err != nil {
			return err
		}
		body = bytes.NewReader(reqBodyBytes)
	}

	resp, err := n.do(uri, method, body, "application/json", nil)
	if err != nil {
		return err
	}

	if response != nil {
		d := json.NewDecoder(resp.Body)
		err := d.Decode(response)
		if err != nil {
			return err
		}
	}

	return nil
}

// Authenticate a user to get an auth token
func (n *PortainerClient) Authenticate() (token string, err error) {
	PrintVerbose("Getting auth token...")

	reqBody := AuthenticateUserRequest{
		Username: n.user,
		Password: n.password,
	}

	respBody := AuthenticateUserResponse{}

	err = n.doJSON("auth", http.MethodPost, &reqBody, &respBody)
	if err != nil {
		return
	}

	token = respBody.Jwt

	return
}

// Get endpoints
func (n *PortainerClient) GetEndpoints() (endpoints []EndpointSubset, err error) {
	PrintVerbose("Getting endpoints...")
	err = n.doJSON("endpoints", http.MethodGet, nil, &endpoints)
	return
}

// Get stacks, optionally filtered by swarmId and endpointId
func (n *PortainerClient) GetStacks(swarmId string, endpointId uint32) (stacks []Stack, err error) {
	PrintVerbose("Getting stacks...")

	filter := StackListFilter{
		SwarmId:    swarmId,
		EndpointId: endpointId,
	}

	filterJsonBytes, _ := json.Marshal(filter)
	filterJsonString := string(filterJsonBytes)

	err = n.doJSON(fmt.Sprintf("stacks?filters=%s", filterJsonString), http.MethodGet, nil, &stacks)
	return
}

// Create swarm stack
func (n *PortainerClient) CreateSwarmStack(stackName string, environmentVariables []StackEnv, stackFileContent string, swarmClusterId string, endpointId string) (err error) {
	PrintVerbose("Deploying stack...")

	reqBody := StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		SwarmID:          swarmClusterId,
		StackFileContent: stackFileContent,
	}

	err = n.doJSON(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%s", 1, "string", endpointId), http.MethodPost, &reqBody, nil)
	return
}

// Create compose stack
func (n *PortainerClient) CreateComposeStack(stackName string, environmentVariables []StackEnv, stackFileContent string, endpointId string) (err error) {
	PrintVerbose("Deploying stack...")

	reqBody := StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
	}

	err = n.doJSON(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%s", 2, "string", endpointId), http.MethodPost, &reqBody, nil)
	return
}

// Update stack
func (n *PortainerClient) UpdateStack(stack Stack, environmentVariables []StackEnv, stackFileContent string, prune bool, endpointId string) (err error) {
	PrintVerbose("Updating stack...")

	reqBody := StackUpdateRequest{
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
		Prune:            prune,
	}

	err = n.doJSON(fmt.Sprintf("stacks/%v?endpointId=%s", stack.Id, endpointId), http.MethodPut, &reqBody, nil)
	return
}

// Delete stack
func (n *PortainerClient) DeleteStack(stackId uint32) (err error) {
	PrintVerbose("Deleting stack...")

	err = n.doJSON(fmt.Sprintf("stacks/%d", stackId), http.MethodDelete, nil, nil)
	return
}

// Get stack file content
func (n *PortainerClient) GetStackFileContent(stackId uint32) (content string, err error) {
	PrintVerbose("Getting stack file content...")

	var respBody StackFileInspectResponse

	err = n.doJSON(fmt.Sprintf("stacks/%v/file", stackId), http.MethodGet, nil, &respBody)
	if err != nil {
		return
	}

	content = respBody.StackFileContent

	return
}

// Get endpoint Docker info
func (n *PortainerClient) GetEndpointDockerInfo(endpointId string) (info map[string]interface{}, err error) {
	PrintVerbose("Getting endpoint Docker info...")

	err = n.doJSON(fmt.Sprintf("endpoints/%v/docker/info", endpointId), http.MethodGet, nil, &info)
	return
}

// Get Portainer status info
func (n *PortainerClient) GetStatus() (status Status, err error) {
	err = n.doJSON("status", http.MethodGet, nil, &status)
	return
}

// Get a clone of the client
func (n *PortainerClient) Clone() (c *PortainerClient, err error) {
	c = &PortainerClient{
		url:           n.url,
		user:          n.user,
		password:      n.password,
		token:         n.token,
		doNotUseToken: n.doNotUseToken,
	}

	c.Timeout = n.Timeout

	c.Transport = n.Transport

	return
}

// Create a new client
func NewClient(config ClientConfig) (c *PortainerClient, err error) {
	apiUrl, err := url.Parse(strings.TrimRight(config.Url, "/") + "/api/")
	if err != nil {
		return
	}

	c = &PortainerClient{
		url:      apiUrl,
		user:     config.User,
		password: config.Password,
		token:    config.Token,
	}

	c.Timeout = config.Timeout

	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure,
		},
	}

	return
}

// Get the cached client or a new one
func GetClient() (c *PortainerClient, err error) {
	if client == nil {
		client, err = NewClient(ClientConfig{
			Url:      viper.GetString("url"),
			User:     viper.GetString("user"),
			Password: viper.GetString("password"),
			Token:    viper.GetString("auth-token"),
		})
	}
	c = client
	return
}

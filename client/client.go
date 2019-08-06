package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type StackListFilter struct {
	SwarmId    string `json:",omitempty"`
	EndpointId uint32 `json:",omitempty"`
}

type Config struct {
	Url           string
	User          string
	Password      string
	Token         string
	DoNotUseToken bool
}

type PortainerClient interface {
	Authenticate() (token string, err error)
	GetEndpoints() ([]EndpointSubset, error)
	GetStacks(swarmId string, endpointId uint32) ([]Stack, error)
	CreateSwarmStack(stackName string, environmentVariables []StackEnv, stackFileContent string, swarmClusterId string, endpointId string) error
	CreateComposeStack(stackName string, environmentVariables []StackEnv, stackFileContent string, endpointId string) error
	UpdateStack(stack Stack, environmentVariables []StackEnv, stackFileContent string, prune bool, endpointId string) error
	DeleteStack(stackId uint32) error
	GetStackFileContent(stackId uint32) (content string, err error)
	GetEndpointDockerInfo(endpointId string) (info map[string]interface{}, err error)
	GetStatus() (Status, error)
	BeforeRequest(hook func(req *http.Request) (err error))
	AfterResponse(hook func(resp *http.Response) (err error))
}

type portainerClientImp struct {
	httpClient         *http.Client
	url                *url.URL
	user               string
	password           string
	token              string
	doNotUseToken      bool
	beforeRequestHooks []func(req *http.Request) (err error)
	afterResponseHooks []func(resp *http.Response) (err error)
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
func (n *portainerClientImp) do(uri, method string, request io.Reader, requestType string, headers http.Header) (resp *http.Response, err error) {
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
			n.token, err = n.Authenticate()
			if err != nil {
				return
			}
		}
		req.Header.Set("Authorization", "Bearer "+n.token)
	}

	// Run all "before request" hooks
	for i := 0; i < len(n.beforeRequestHooks); i++ {
		err = n.beforeRequestHooks[i](req)
		if err != nil {
			return
		}
	}

	resp, err = n.httpClient.Do(req)
	if err != nil {
		return
	}

	// Run all "after response" hooks
	for i := 0; i < len(n.afterResponseHooks); i++ {
		err = n.afterResponseHooks[i](resp)
		if err != nil {
			return
		}
	}

	err = checkResponseForErrors(resp)
	if err != nil {
		return
	}

	return
}

// Do a JSON http request
func (n *portainerClientImp) doJSON(uri, method string, request interface{}, response interface{}) error {
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

func (n *portainerClientImp) BeforeRequest(hook func(req *http.Request) (err error)) {
	n.beforeRequestHooks = append(n.beforeRequestHooks, hook)
}

func (n *portainerClientImp) AfterResponse(hook func(resp *http.Response) (err error)) {
	n.afterResponseHooks = append(n.afterResponseHooks, hook)
}

// Authenticate a user to get an auth token
func (n *portainerClientImp) Authenticate() (token string, err error) {
	reqBody := AuthenticateUserRequest{
		Username: n.user,
		Password: n.password,
	}

	respBody := AuthenticateUserResponse{}

	previousDoNotUseTokenValue := n.doNotUseToken
	n.doNotUseToken = true

	err = n.doJSON("auth", http.MethodPost, &reqBody, &respBody)
	if err != nil {
		return
	}

	n.doNotUseToken = previousDoNotUseTokenValue

	token = respBody.Jwt

	return
}

// Get endpoints
func (n *portainerClientImp) GetEndpoints() (endpoints []EndpointSubset, err error) {
	err = n.doJSON("endpoints", http.MethodGet, nil, &endpoints)
	return
}

// Get stacks, optionally filtered by swarmId and endpointId
func (n *portainerClientImp) GetStacks(swarmId string, endpointId uint32) (stacks []Stack, err error) {
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
func (n *portainerClientImp) CreateSwarmStack(stackName string, environmentVariables []StackEnv, stackFileContent string, swarmClusterId string, endpointId string) (err error) {
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
func (n *portainerClientImp) CreateComposeStack(stackName string, environmentVariables []StackEnv, stackFileContent string, endpointId string) (err error) {
	reqBody := StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
	}

	err = n.doJSON(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%s", 2, "string", endpointId), http.MethodPost, &reqBody, nil)
	return
}

// Update stack
func (n *portainerClientImp) UpdateStack(stack Stack, environmentVariables []StackEnv, stackFileContent string, prune bool, endpointId string) (err error) {
	reqBody := StackUpdateRequest{
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
		Prune:            prune,
	}

	err = n.doJSON(fmt.Sprintf("stacks/%v?endpointId=%s", stack.Id, endpointId), http.MethodPut, &reqBody, nil)
	return
}

// Delete stack
func (n *portainerClientImp) DeleteStack(stackId uint32) (err error) {
	err = n.doJSON(fmt.Sprintf("stacks/%d", stackId), http.MethodDelete, nil, nil)
	return
}

// Get stack file content
func (n *portainerClientImp) GetStackFileContent(stackId uint32) (content string, err error) {
	var respBody StackFileInspectResponse

	err = n.doJSON(fmt.Sprintf("stacks/%v/file", stackId), http.MethodGet, nil, &respBody)
	if err != nil {
		return
	}

	content = respBody.StackFileContent

	return
}

// Get endpoint Docker info
func (n *portainerClientImp) GetEndpointDockerInfo(endpointId string) (info map[string]interface{}, err error) {
	err = n.doJSON(fmt.Sprintf("endpoints/%v/docker/info", endpointId), http.MethodGet, nil, &info)
	return
}

// Get Portainer status info
func (n *portainerClientImp) GetStatus() (status Status, err error) {
	err = n.doJSON("status", http.MethodGet, nil, &status)
	return
}

// Create a new client
func NewClient(httpClient *http.Client, config Config) (c PortainerClient, err error) {
	apiUrl, err := url.Parse(strings.TrimRight(config.Url, "/") + "/api/")
	if err != nil {
		return
	}

	c = &portainerClientImp{
		httpClient: httpClient,
		url:        apiUrl,
		user:       config.User,
		password:   config.Password,
		token:      config.Token,
	}

	return
}

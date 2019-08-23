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

	portainer "github.com/portainer/portainer/api"
)

type StackListFilter struct {
	SwarmID    string               `json:"SwarmId,omitempty"`
	EndpointID portainer.EndpointID `json:"EndpointId,omitempty"`
}

type Config struct {
	URL           *url.URL
	User          string
	Password      string
	Token         string
	UserAgent     string
	DoNotUseToken bool
}

type PortainerClient interface {
	// Authenticate a user to get an auth token
	Authenticate() (token string, err error)

	// Get endpoints
	GetEndpoints() ([]portainer.Endpoint, error)

	// Get endpoint groups
	GetEndpointGroups() ([]portainer.EndpointGroup, error)

	// Get stacks, optionally filtered by swarmId and endpointId
	GetStacks(swarmID string, endpointID portainer.EndpointID) ([]portainer.Stack, error)

	// Create swarm stack
	CreateSwarmStack(stackName string, environmentVariables []portainer.Pair, stackFileContent string, swarmClusterID string, endpointID portainer.EndpointID) (stack portainer.Stack, err error)

	// Create compose stack
	CreateComposeStack(stackName string, environmentVariables []portainer.Pair, stackFileContent string, endpointID portainer.EndpointID) (stack portainer.Stack, err error)

	// Update stack
	UpdateStack(stack portainer.Stack, environmentVariables []portainer.Pair, stackFileContent string, prune bool, endpointID portainer.EndpointID) error

	// Delete stack
	DeleteStack(stackID portainer.StackID) error

	// Get stack file content
	GetStackFileContent(stackID portainer.StackID) (content string, err error)

	// Get endpoint Docker info
	GetEndpointDockerInfo(endpointID portainer.EndpointID) (info map[string]interface{}, err error)

	// Get Portainer status info
	GetStatus() (portainer.Status, error)

	// Run a function before sending a request to Portainer
	BeforeRequest(hook func(req *http.Request) (err error))

	// Run a function after receiving a response from Portainer
	AfterResponse(hook func(resp *http.Response) (err error))
}

type portainerClientImp struct {
	httpClient         *http.Client
	url                *url.URL
	user               string
	password           string
	token              string
	userAgent          string
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
	requestURL, err := n.url.Parse(uri)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, requestURL.String(), request)
	if err != nil {
		return
	}

	if headers != nil {
		req.Header = headers
	}

	if request != nil {
		req.Header.Set("Content-Type", requestType)
		req.Header.Set("User-Agent", n.userAgent)
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

func (n *portainerClientImp) GetEndpoints() (endpoints []portainer.Endpoint, err error) {
	err = n.doJSON("endpoints", http.MethodGet, nil, &endpoints)
	return
}

func (n *portainerClientImp) GetEndpointGroups() (endpointGroups []portainer.EndpointGroup, err error) {
	err = n.doJSON("endpoint_groups", http.MethodGet, nil, &endpointGroups)
	return
}

func (n *portainerClientImp) GetStacks(swarmID string, endpointID portainer.EndpointID) (stacks []portainer.Stack, err error) {
	filter := StackListFilter{
		SwarmID:    swarmID,
		EndpointID: endpointID,
	}

	filterJSONBytes, _ := json.Marshal(filter)
	filterJSONString := string(filterJSONBytes)

	err = n.doJSON(fmt.Sprintf("stacks?filters=%s", filterJSONString), http.MethodGet, nil, &stacks)
	return
}

func (n *portainerClientImp) CreateSwarmStack(stackName string, environmentVariables []portainer.Pair, stackFileContent string, swarmClusterID string, endpointID portainer.EndpointID) (stack portainer.Stack, err error) {
	reqBody := StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		SwarmID:          swarmClusterID,
		StackFileContent: stackFileContent,
	}

	err = n.doJSON(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%v", 1, "string", endpointID), http.MethodPost, &reqBody, &stack)
	return
}

func (n *portainerClientImp) CreateComposeStack(stackName string, environmentVariables []portainer.Pair, stackFileContent string, endpointID portainer.EndpointID) (stack portainer.Stack, err error) {
	reqBody := StackCreateRequest{
		Name:             stackName,
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
	}

	err = n.doJSON(fmt.Sprintf("stacks?type=%v&method=%s&endpointId=%v", 2, "string", endpointID), http.MethodPost, &reqBody, &stack)
	return
}

func (n *portainerClientImp) UpdateStack(stack portainer.Stack, environmentVariables []portainer.Pair, stackFileContent string, prune bool, endpointID portainer.EndpointID) (err error) {
	reqBody := StackUpdateRequest{
		Env:              environmentVariables,
		StackFileContent: stackFileContent,
		Prune:            prune,
	}

	err = n.doJSON(fmt.Sprintf("stacks/%v?endpointId=%v", stack.ID, endpointID), http.MethodPut, &reqBody, nil)
	return
}

func (n *portainerClientImp) DeleteStack(stackID portainer.StackID) (err error) {
	err = n.doJSON(fmt.Sprintf("stacks/%d", stackID), http.MethodDelete, nil, nil)
	return
}

func (n *portainerClientImp) GetStackFileContent(stackID portainer.StackID) (content string, err error) {
	var respBody StackFileInspectResponse

	err = n.doJSON(fmt.Sprintf("stacks/%v/file", stackID), http.MethodGet, nil, &respBody)
	if err != nil {
		return
	}

	content = respBody.StackFileContent

	return
}

func (n *portainerClientImp) GetEndpointDockerInfo(endpointID portainer.EndpointID) (info map[string]interface{}, err error) {
	err = n.doJSON(fmt.Sprintf("endpoints/%v/docker/info", endpointID), http.MethodGet, nil, &info)
	return
}

func (n *portainerClientImp) GetStatus() (status portainer.Status, err error) {
	err = n.doJSON("status", http.MethodGet, nil, &status)
	return
}

// Create a new client
func NewClient(httpClient *http.Client, config Config) PortainerClient {
	return &portainerClientImp{
		httpClient: httpClient,
		url:        config.URL,
		user:       config.User,
		password:   config.Password,
		token:      config.Token,
		userAgent:  config.UserAgent,
	}
}

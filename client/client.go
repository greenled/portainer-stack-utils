package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	portainer "github.com/portainer/portainer/api"
)

// Config represents a Portainer client configuration
type Config struct {
	URL           *url.URL
	User          string
	Password      string
	Token         string
	UserAgent     string
	DoNotUseToken bool
}

// PortainerClient represents a Portainer API client
type PortainerClient interface {
	// AuthenticateUser a user to get an auth token
	AuthenticateUser(options AuthenticateUserOptions) (token string, err error)

	// Get endpoints
	EndpointList() ([]portainer.Endpoint, error)

	// Get endpoint groups
	EndpointGroupList() ([]portainer.EndpointGroup, error)

	// Get stacks, optionally filtered by swarmId and endpointId
	StackList(options StackListOptions) ([]portainer.Stack, error)

	// Create swarm stack
	StackCreateSwarm(options StackCreateSwarmOptions) (stack portainer.Stack, err error)

	// Create compose stack
	StackCreateCompose(options StackCreateComposeOptions) (stack portainer.Stack, err error)

	// Update stack
	StackUpdate(options StackUpdateOptions) error

	// Delete stack
	StackDelete(stackID portainer.StackID) error

	// Get stack file content
	StackFileInspect(stackID portainer.StackID) (content string, err error)

	// Get endpoint Docker info
	EndpointDockerInfo(endpointID portainer.EndpointID) (info map[string]interface{}, err error)

	// Get Portainer status info
	Status() (portainer.Status, error)

	// Run a function before sending a request to Portainer
	BeforeRequest(hook func(req *http.Request) (err error))

	// Run a function after receiving a response from Portainer
	AfterResponse(hook func(resp *http.Response) (err error))

	// Proxy proxies a request to /endpoint/{id}/docker and returns its result
	Proxy(endpointID portainer.EndpointID, req *http.Request) (resp *http.Response, err error)
}

type portainerClientImp struct {
	httpClient         *http.Client
	url                *url.URL
	user               string
	password           string
	token              string
	userAgent          string
	beforeRequestHooks []func(req *http.Request) (err error)
	afterResponseHooks []func(resp *http.Response) (err error)
}

// Do an http request
func (n *portainerClientImp) do(uri, method string, requestBody io.Reader, headers http.Header) (resp *http.Response, err error) {
	requestURL, err := n.url.Parse(uri)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, requestURL.String(), requestBody)
	if err != nil {
		return
	}

	if headers != nil {
		req.Header = headers
	}

	// Set user agent header
	req.Header.Set("User-Agent", n.userAgent)

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

func (n *portainerClientImp) doWithToken(uri, method string, requestBody io.Reader, headers http.Header) (resp *http.Response, err error) {
	// Ensure there is an auth token
	if n.token == "" {
		n.token, err = n.AuthenticateUser(AuthenticateUserOptions{
			Username: n.user,
			Password: n.password,
		})
		if err != nil {
			return
		}
	}
	headers.Set("Authorization", "Bearer "+n.token)

	return n.do(uri, method, requestBody, headers)
}

// Do a JSON http request
func (n *portainerClientImp) doJSON(uri, method string, headers http.Header, requestBody interface{}, responseBody interface{}) error {
	// Encode request body, if any
	var body io.Reader
	if requestBody != nil {
		reqBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(reqBodyBytes)
	}

	// Set content type header
	headers.Set("Content-Type", "application/json")

	resp, err := n.do(uri, method, body, headers)
	if err != nil {
		return err
	}

	// Decode response body, if any
	if responseBody != nil {
		d := json.NewDecoder(resp.Body)
		err := d.Decode(responseBody)
		if err != nil {
			return err
		}
	}

	return nil
}

// Do a JSON http request with an auth token
func (n *portainerClientImp) doJSONWithToken(uri, method string, headers http.Header, request interface{}, response interface{}) (err error) {
	// Ensure there is an auth token
	if n.token == "" {
		n.token, err = n.AuthenticateUser(AuthenticateUserOptions{
			Username: n.user,
			Password: n.password,
		})
		if err != nil {
			return
		}
	}
	headers.Set("Authorization", "Bearer "+n.token)

	return n.doJSON(uri, method, headers, request, response)
}

func (n *portainerClientImp) BeforeRequest(hook func(req *http.Request) (err error)) {
	n.beforeRequestHooks = append(n.beforeRequestHooks, hook)
}

func (n *portainerClientImp) AfterResponse(hook func(resp *http.Response) (err error)) {
	n.afterResponseHooks = append(n.afterResponseHooks, hook)
}

// NewClient creates a new Portainer API client
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

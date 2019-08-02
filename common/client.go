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
	"time"

	"github.com/spf13/viper"
)

var client *PortainerClient

type PortainerClient struct {
	http.Client
	url   *url.URL
	token string
}

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

	if n.token != "" {
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

func (n *PortainerClient) DoJSON(uri, method string, request interface{}, response interface{}) error {
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

func (n *PortainerClient) Authenticate(user, password string) (token string, err error) {
	PrintVerbose("Getting auth token...")

	reqBody := AuthenticateUserRequest{
		Username: viper.GetString("user"),
		Password: viper.GetString("password"),
	}

	respBody := AuthenticateUserResponse{}

	err = n.DoJSON("auth", http.MethodPost, &reqBody, &respBody)
	if err != nil {
		return
	}

	token = respBody.Jwt

	return
}

type clientConfig struct {
	Url      string
	User     string
	Password string
	Token    string
	Insecure bool
	Timeout  time.Duration
}

func newClient(config clientConfig) (c *PortainerClient, err error) {
	apiUrl, err := url.Parse(config.Url + "/api/")
	if err != nil {
		return
	}

	c = &PortainerClient{
		url: apiUrl,
	}

	c.Timeout = config.Timeout

	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure,
		},
	}

	if config.Token != "" {
		c.token = config.Token
	} else {
		c.token, err = c.Authenticate(config.User, config.Password)
		if err != nil {
			return nil, err
		}
		PrintDebug(fmt.Sprintf("Auth token: %s", c.token))
	}

	return
}

func GetClient() (c *PortainerClient, err error) {
	if client == nil {
		client, err = newClient(clientConfig{
			Url:      viper.GetString("url"),
			User:     viper.GetString("user"),
			Password: viper.GetString("password"),
			Token:    viper.GetString("auth-token"),
		})
	}
	c = client
	return
}

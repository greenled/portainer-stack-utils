package common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/greenled/portainer-stack-utils/util"

	"github.com/greenled/portainer-stack-utils/client"

	"github.com/spf13/viper"
)

var cachedClient client.PortainerClient

// Get the cached client or a new one
func GetClient() (c client.PortainerClient, err error) {
	if cachedClient == nil {
		cachedClient, err = GetDefaultClient()
		if err != nil {
			return
		}
	}
	return cachedClient, nil
}

// Get the default client
func GetDefaultClient() (c client.PortainerClient, err error) {
	c, err = client.NewClient(GetDefaultHttpClient(), GetDefaultClientConfig())
	if err != nil {
		return
	}

	c.BeforeRequest(func(req *http.Request) (err error) {
		var bodyString string
		if req.Body != nil {
			bodyBytes, readErr := ioutil.ReadAll(req.Body)
			defer req.Body.Close()
			if readErr != nil {
				return readErr
			}
			bodyString = string(bodyBytes)
			req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		}

		util.PrintDebug(fmt.Sprintf(`Request
---
Method: %s
URL: %s
Body:
%s
---`, req.Method, req.URL.String(), string(bodyString)))

		return
	})

	c.AfterResponse(func(resp *http.Response) (err error) {
		var bodyString string
		if resp.Body != nil {
			bodyBytes, readErr := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if readErr != nil {
				return readErr
			}
			bodyString = string(bodyBytes)
			resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		}

		util.PrintDebug(fmt.Sprintf(`Response
---
Status: %s
Body:
%s
---`, resp.Status, bodyString))

		return
	})

	return
}

// Get the default config for a client
func GetDefaultClientConfig() client.Config {
	return client.Config{
		Url:           viper.GetString("url"),
		User:          viper.GetString("user"),
		Password:      viper.GetString("password"),
		Token:         viper.GetString("auth-token"),
		DoNotUseToken: false,
	}
}

// Get the default http client for a Portainer client
func GetDefaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: viper.GetDuration("timeout"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: viper.GetBool("insecure"),
			},
		},
	}
}

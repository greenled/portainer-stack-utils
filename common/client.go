package common

import (
	"crypto/tls"
	"net/http"

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
	return client.NewClient(GetDefaultHttpClient(), GetDefaultClientConfig())
}

// Get the default config for a client
func GetDefaultClientConfig() client.ClientConfig {
	return client.ClientConfig{
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

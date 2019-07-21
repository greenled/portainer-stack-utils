package common

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"net/http"
)

func NewHttpClient() http.Client {
	// Create HTTP transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: viper.GetBool("insecure"),
		},
	}

	// Create HTTP client
	return http.Client{
		Transport: tr,
		Timeout:   viper.GetDuration("timeout"),
	}
}

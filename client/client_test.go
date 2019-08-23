package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readRequestBodyAsJSON(req *http.Request, body *map[string]interface{}) (err error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	err = json.Unmarshal(bodyBytes, body)
	return
}

func writeResponseBodyAsJSON(w http.ResponseWriter, body map[string]interface{}) (err error) {
	bodyBytes, err := json.Marshal(body)
	fmt.Fprintln(w, string(bodyBytes))
	return
}

func TestNewClient(t *testing.T) {
	apiURL, _ := url.Parse("http://validurl.com/api")

	validClient := NewClient(http.DefaultClient, Config{
		URL: apiURL,
	})
	assert.NotNil(t, validClient)
}

func TestClientAuthenticates(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body map[string]interface{}
		err := readRequestBodyAsJSON(req, &body)

		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.RequestURI, "/api/auth")
		assert.NotNil(t, req.Header["Content-Type"])
		assert.NotNil(t, req.Header["Content-Type"][0])
		assert.Equal(t, req.Header["Content-Type"][0], "application/json")
		assert.NotNil(t, req.Header["User-Agent"])
		assert.NotNil(t, req.Header["User-Agent"][0])
		assert.Equal(t, req.Header["User-Agent"][0], "GE007")
		assert.Nil(t, err)
		assert.NotNil(t, body["Username"])
		assert.Equal(t, body["Username"], "admin")
		assert.NotNil(t, body["Password"])
		assert.Equal(t, body["Password"], "a")

		writeResponseBodyAsJSON(w, map[string]interface{}{
			"jwt": "somerandomtoken",
		})
	}))
	defer ts.Close()

	apiURL, _ := url.Parse(ts.URL + "/api/")

	customClient := NewClient(ts.Client(), Config{
		URL:       apiURL,
		User:      "admin",
		Password:  "a",
		UserAgent: "GE007",
	})
	token, err := customClient.Authenticate()
	assert.Nil(t, err)
	assert.Equal(t, token, "somerandomtoken")
}

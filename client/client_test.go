package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readRequestBodyAsJson(req *http.Request, body *map[string]interface{}) (err error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	err = json.Unmarshal(bodyBytes, body)
	return
}

func writeResponseBodyAsJson(w http.ResponseWriter, body map[string]interface{}) (err error) {
	bodyBytes, err := json.Marshal(body)
	fmt.Fprintln(w, string(bodyBytes))
	return
}

func TestNewClient(t *testing.T) {
	validClient, err := NewClient(http.DefaultClient, Config{
		Url: "http://validurl.com",
	})
	assert.NotNil(t, validClient)
	assert.Nil(t, err)
}

func TestClientAuthenticates(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body map[string]interface{}
		err := readRequestBodyAsJson(req, &body)

		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.RequestURI, "/api/auth")
		assert.NotNil(t, req.Header["Content-Type"])
		assert.NotNil(t, req.Header["Content-Type"][0])
		assert.Equal(t, req.Header["Content-Type"][0], "application/json")
		assert.Nil(t, err)
		assert.NotNil(t, body["Username"])
		assert.Equal(t, body["Username"], "admin")
		assert.NotNil(t, body["Password"])
		assert.Equal(t, body["Password"], "a")

		writeResponseBodyAsJson(w, map[string]interface{}{
			"jwt": "somerandomtoken",
		})
	}))
	defer ts.Close()

	customClient, err := NewClient(ts.Client(), Config{
		Url:      ts.URL,
		User:     "admin",
		Password: "a",
	})
	assert.Nil(t, err)
	token, err := customClient.Authenticate()
	assert.Nil(t, err)
	assert.Equal(t, token, "somerandomtoken")
}

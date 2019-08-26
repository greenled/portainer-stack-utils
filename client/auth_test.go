package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	token, err := customClient.Auth()
	assert.Nil(t, err)
	assert.Equal(t, token, "somerandomtoken")
}

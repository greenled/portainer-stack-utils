package client

import (
	"encoding/json"
	"fmt"
	"io"
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

func Test_portainerClientImp_do(t *testing.T) {
	type fields struct {
		user               string
		password           string
		token              string
		userAgent          string
		beforeRequestHooks []func(req *http.Request) (err error)
		afterResponseHooks []func(resp *http.Response) (err error)
		server             *httptest.Server
		beforeFunctionCall func(t *testing.T, tt *fields)
	}
	type args struct {
		uri         string
		method      string
		requestBody io.Reader
		headers     http.Header
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantRespCheck func(resp *http.Response) bool
		wantErr       bool
	}{
		{
			name: "error on bad URI",
			fields: fields{
				server: httptest.NewUnstartedServer(nil),
			},
			args: args{
				uri: string(0x7f),
			},
			wantErr: true,
		},
		{
			name: "error on bad method",
			fields: fields{
				server: httptest.NewUnstartedServer(nil),
			},
			args: args{
				method: "WOLOLO?",
			},
			wantErr: true,
		},
		{
			name: "extra headers are added",
			fields: fields{
				server: httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.Header.Get("Some-Header"), "value")
				})),
			},
			args: args{
				headers: http.Header{
					"Some-Header": []string{
						"value",
					},
				},
			},
		},
		{
			name: "returns error on http error",
			fields: fields{
				token:  "token",
				server: httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})),
				beforeFunctionCall: func(t *testing.T, tt *fields) {
					tt.server.Close()
				},
			},
			wantErr: true,
		},
		{
			name: "returns error on response error",
			fields: fields{
				server: httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.server.Start()
			defer tt.fields.server.Close()

			apiURL, _ := url.Parse(tt.fields.server.URL + "/api/")

			n := &portainerClientImp{
				httpClient:         tt.fields.server.Client(),
				url:                apiURL,
				user:               tt.fields.user,
				password:           tt.fields.password,
				token:              tt.fields.token,
				userAgent:          tt.fields.userAgent,
				beforeRequestHooks: tt.fields.beforeRequestHooks,
				afterResponseHooks: tt.fields.afterResponseHooks,
			}

			if tt.fields.beforeFunctionCall != nil {
				tt.fields.beforeFunctionCall(t, &tt.fields)
			}
			gotResp, err := n.do(tt.args.uri, tt.args.method, tt.args.requestBody, tt.args.headers)

			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantRespCheck != nil {
				assert.True(t, tt.wantRespCheck(gotResp))
			}
		})
	}
}

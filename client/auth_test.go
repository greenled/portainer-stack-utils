package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_portainerClientImp_AuthenticateUser(t *testing.T) {
	type fields struct {
		server *httptest.Server
	}
	type args struct {
		options AuthenticateUserOptions
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			name: "valid username and password authenticates (happy path)",
			fields: fields{
				server: httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.Method, http.MethodPost)
					assert.Equal(t, req.RequestURI, "/api/auth")

					var body map[string]interface{}
					err := readRequestBodyAsJSON(req, &body)
					assert.Nil(t, err)

					assert.NotNil(t, body["Username"])
					assert.Equal(t, body["Username"], "admin")
					assert.NotNil(t, body["Password"])
					assert.Equal(t, body["Password"], "a")

					writeResponseBodyAsJSON(w, map[string]interface{}{
						"jwt": "token",
					})
				})),
			},
			args: args{
				options: AuthenticateUserOptions{
					Username: "admin",
					Password: "a",
				},
			},
			wantToken: "token",
		},
		{
			name: "invalid username and password does not authenticate",
			fields: fields{
				server: httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusUnprocessableEntity)
					writeResponseBodyAsJSON(w, map[string]interface{}{
						"Err":     "Invalid credentials",
						"Details": "Unauthorized",
					})
				})),
			},
			args: args{
				options: AuthenticateUserOptions{
					Username: "admin",
					Password: "a",
				},
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
				httpClient: tt.fields.server.Client(),
				url:        apiURL,
			}

			gotToken, err := n.AuthenticateUser(tt.args.options)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantToken, gotToken)
		})
	}
}

package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	portainer "github.com/portainer/portainer/api"
	"github.com/stretchr/testify/assert"
)

func TestGetTranslatedStackType(t *testing.T) {
	type args struct {
		t portainer.StackType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "swarm stack type",
			args: args{
				t: portainer.DockerSwarmStack,
			},
			want: "swarm",
		},
		{
			name: "compose stack type",
			args: args{
				t: portainer.DockerComposeStack,
			},
			want: "compose",
		},
		{
			name: "unknown stack type",
			args: args{
				t: 100,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetTranslatedStackType(tt.args.t))
		})
	}
}

func Test_checkResponseForErrors(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "generic error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusNotFound,
					}
					bodyBytes, _ := json.Marshal(map[string]interface{}{
						"Err":     "Error",
						"Details": "Not found",
					})
					resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
					return
				}(),
			},
			wantErr: true,
		},
		{
			name: "non generic error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("Err"))),
					}
					return
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, checkResponseForErrors(tt.args.resp) != nil)
		})
	}
}

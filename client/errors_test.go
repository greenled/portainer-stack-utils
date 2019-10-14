package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericError_Error(t *testing.T) {
	type fields struct {
		Err     string
		Details string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "error with message and details",
			fields: fields{
				Err:     "error",
				Details: "details",
			},
			want: "error: details",
		},
		{
			name: "error with message and no details",
			fields: fields{
				Err: "error",
			},
			want: "error",
		},
		{
			name: "error with no error message and details",
			fields: fields{
				Details: "details",
			},
			want: ": details",
		},
		{
			name:   "error with no error message and no details",
			fields: fields{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &GenericError{
				Err:     tt.fields.Err,
				Details: tt.fields.Details,
			}
			assert.Equal(t, tt.want, e.Error())
		})
	}
}

func Test_getResponseHTTPError(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "bad request (generic) error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusBadRequest,
					}
					bodyBytes, _ := json.Marshal(map[string]interface{}{
						"Err":     "Error",
						"Details": "Bad request",
					})
					resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
					return
				}(),
			},
			wantErr: &GenericError{
				Code:    http.StatusBadRequest,
				Err:     "Error",
				Details: "Bad request",
			},
		},
		{
			name: "forbidden (generic) error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusForbidden,
					}
					bodyBytes, _ := json.Marshal(map[string]interface{}{
						"Err":     "Error",
						"Details": "Forbidden",
					})
					resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
					return
				}(),
			},
			wantErr: &GenericError{
				Code:    http.StatusForbidden,
				Err:     "Error",
				Details: "Forbidden",
			},
		},
		{
			name: "not found (generic) error",
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
			wantErr: &GenericError{
				Code:    http.StatusNotFound,
				Err:     "Error",
				Details: "Not found",
			},
		},
		{
			name: "conflict (generic) error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusConflict,
					}
					bodyBytes, _ := json.Marshal(map[string]interface{}{
						"Err":     "Error",
						"Details": "Conflict",
					})
					resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
					return
				}(),
			},
			wantErr: &GenericError{
				Code:    http.StatusConflict,
				Err:     "Error",
				Details: "Conflict",
			},
		},
		{
			name: "internal server error (generic) error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusInternalServerError,
					}
					bodyBytes, _ := json.Marshal(map[string]interface{}{
						"Err":     "Error",
						"Details": "Internal server error",
					})
					resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
					return
				}(),
			},
			wantErr: &GenericError{
				Code:    http.StatusInternalServerError,
				Err:     "Error",
				Details: "Internal server error",
			},
		},
		{
			name: "service unavailable (generic) error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusServiceUnavailable,
					}
					bodyBytes, _ := json.Marshal(map[string]interface{}{
						"Err":     "Error",
						"Details": "Service unavailable",
					})
					resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
					return
				}(),
			},
			wantErr: &GenericError{
				Code:    http.StatusServiceUnavailable,
				Err:     "Error",
				Details: "Service unavailable",
			},
		},
		{
			name: "method not allowed (non generic) error",
			args: args{
				resp: func() (resp *http.Response) {
					resp = &http.Response{
						StatusCode: http.StatusMethodNotAllowed,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("Err"))),
					}
					return
				}(),
			},
			wantErr: errors.New("Err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, getResponseHTTPError(tt.args.resp))
		})
	}
}

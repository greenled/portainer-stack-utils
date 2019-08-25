package client

import (
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

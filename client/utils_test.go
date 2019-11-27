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

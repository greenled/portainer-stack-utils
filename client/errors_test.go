package client

import (
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

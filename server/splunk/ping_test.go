package splunk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splunk_Ping(t *testing.T) {
	is := assert.New(t)

	type fields struct {
		Config Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ping success",
			fields: fields{
				Config: Config{
					SplunkUserInfo: User{
						ServerBaseURL: "https://207.154.235.95",
						UserName:      "splunk_admin",
						Password:      "splunk_admin",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ping fail",
			fields: fields{
				Config: Config{
					SplunkUserInfo: User{
						ServerBaseURL: "https://207.154.235.95",
						UserName:      "splunk_admin",
						Password:      "123",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.Config)
			err := s.Ping()
			is.Equal(tt.wantErr, err != nil)
		})
	}
}

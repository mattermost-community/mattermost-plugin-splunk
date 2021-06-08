package splunk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splunk_Logs(t *testing.T) {
	is := assert.New(t)

	type fields struct {
		Config Config
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				Config: Config{
					Dependencies: nil,
					SplunkUserInfo: User{
						ServerBaseURL: "https://207.154.235.95",
						UserName:      "splunk_admin",
						Password:      "splunk_admin",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.Config)
			logs, err := s.Logs("downloaded-logs-20210119-171858.csv")
			is.NoError(err)
			fmt.Println(logs)
		})
	}
}

func Test_splunk_ListLogs(t *testing.T) {
	type fields struct {
		Config Config
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				Config: Config{
					Dependencies: nil,
					SplunkUserInfo: User{
						ServerBaseURL: "https://207.154.235.95",
						UserName:      "splunk_admin",
						Password:      "splunk_admin",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.Config)
			fmt.Println(s.ListLogs())
		})
	}
}

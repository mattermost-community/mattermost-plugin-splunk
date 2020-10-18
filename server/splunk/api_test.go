package splunk

import (
	"testing"
)

func Test_splunk_Logs(t *testing.T) {
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
					Dependencies:        nil,
					SplunkServerBaseURL: "https://207.154.235.95",
					SplunkUserName:      "bakurits",
					SplunkPassword:      "matarebeli",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.Config)
			s.Logs()
		})
	}
}

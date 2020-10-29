package splunk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
			logs, err := s.Logs("gae_app_module_id_default__logs__2020-10-14T12-56.csv")
			is.NoError(err)
			fmt.Println(logs)
		})
	}
}

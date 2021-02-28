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
						Token:         "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJzcGx1bmtfYWRtaW4gZnJvbSB1YnVudHUtcy0xdmNwdS0xZ2ItZnJhMS0wMSIsInN1YiI6InNwbHVua19hZG1pbiIsImF1ZCI6ImdnIiwiaWRwIjoiU3BsdW5rIiwianRpIjoiMzIzOGFhNDA4MDkxNTI5MDkzMDZhNzYxYTk5MWQ5YjEzZjZkNmE5YmI1ZmMzZGM0NTA5MzhmNjY2MDgyODY1NSIsImlhdCI6MTYxNDUwODYyNywiZXhwIjoxNjE3MTAwNjI3LCJuYnIiOjE2MTQ1MDg2Mjd9.UPOpCe3zsi_dZ3P6GomfRGklVL-ef8DyMXH0DAUPMzp3xAUKp_EFxRbguslCbJ0dU1e6O_DpXzzINSaEKlsWqw",
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
						Token:         "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJzcGx1bmtfYWRtaW4gZnJvbSB1YnVudHUtcy0xdmNwdS0xZ2ItZnJhMS0wMSIsInN1YiI6InNwbHVua19hZG1pbiIsImF1ZCI6ImdnIiwiaWRwIjoiU3BsdW5rIiwianRpIjoiMzIzOGFhNDA4MDkxNTI5MDkzMDZhNzYxYTk5MWQ5YjEzZjZkNmE5YmI1ZmMzZGM0NTA5MzhmNjY2MDgyODY1NSIsImlhdCI6MTYxNDUwODYyNywiZXhwIjoxNjE3MTAwNjI3LCJuYnIiOjE2MTQ1MDg2Mjd9.UPOpCe3zsi_dZ3P6GomfRGklVL-ef8DyMXH0DAUPMzp3xAUKp_EFxRbguslCbJ0dU1e6O_DpXzzINSaEKlsWqw",
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

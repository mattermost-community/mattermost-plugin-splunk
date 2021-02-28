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
						ServerBaseURL: "https://207.154.235.95:8089",
						Token:         "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJzcGx1bmtfYWRtaW4gZnJvbSB1YnVudHUtcy0xdmNwdS0xZ2ItZnJhMS0wMSIsInN1YiI6InNwbHVua19hZG1pbiIsImF1ZCI6ImdnIiwiaWRwIjoiU3BsdW5rIiwianRpIjoiMzIzOGFhNDA4MDkxNTI5MDkzMDZhNzYxYTk5MWQ5YjEzZjZkNmE5YmI1ZmMzZGM0NTA5MzhmNjY2MDgyODY1NSIsImlhdCI6MTYxNDUwODYyNywiZXhwIjoxNjE3MTAwNjI3LCJuYnIiOjE2MTQ1MDg2Mjd9.UPOpCe3zsi_dZ3P6GomfRGklVL-ef8DyMXH0DAUPMzp3xAUKp_EFxRbguslCbJ0dU1e6O_DpXzzINSaEKlsWqw",
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
						ServerBaseURL: "https://207.154.235.95:8089",
						Token:         "bad creds",
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
			if tt.wantErr && err == nil {
				is.Fail("expected error but didn't get any")
			}
			if !tt.wantErr && err != nil {
				is.Failf("", "didn't expect any error but got %v", err)
			}
		})
	}
}

package splunk

import (
	"fmt"
	"testing"

	"github.com/bakurits/mattermost-plugin-splunk/server/store"
	"github.com/bakurits/mattermost-plugin-splunk/server/store/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_splunk_Logs(t *testing.T) {
	ctrl := gomock.NewController(t)
	is := assert.New(t)
	m := mock.NewMockStore(ctrl)

	m.EXPECT().ChangeCurrentUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	defer ctrl.Finish()

	type fields struct {
		User store.SplunkUser
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				User: store.SplunkUser{
					Server: "https://207.154.235.95:8089",
					Token:  "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJzcGx1bmtfYWRtaW4gZnJvbSB1YnVudHUtcy0xdmNwdS0xZ2ItZnJhMS0wMSIsInN1YiI6InNwbHVua19hZG1pbiIsImF1ZCI6ImdnIiwiaWRwIjoiU3BsdW5rIiwianRpIjoiMzIzOGFhNDA4MDkxNTI5MDkzMDZhNzYxYTk5MWQ5YjEzZjZkNmE5YmI1ZmMzZGM0NTA5MzhmNjY2MDgyODY1NSIsImlhdCI6MTYxNDUwODYyNywiZXhwIjoxNjE3MTAwNjI3LCJuYnIiOjE2MTQ1MDg2Mjd9.UPOpCe3zsi_dZ3P6GomfRGklVL-ef8DyMXH0DAUPMzp3xAUKp_EFxRbguslCbJ0dU1e6O_DpXzzINSaEKlsWqw",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSplunk(nil, m)
			s.currentUser = tt.fields.User

			logs, err := s.Logs("downloaded-logs-20210119-171858.csv")
			is.NoError(err)
			fmt.Println(logs)
		})
	}
}

func Test_splunk_ListLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock.NewMockStore(ctrl)

	m.EXPECT().ChangeCurrentUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	defer ctrl.Finish()

	type fields struct {
		User store.SplunkUser
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				User: store.SplunkUser{
					Server: "https://207.154.235.95:8089",
					Token:  "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJzcGx1bmtfYWRtaW4gZnJvbSB1YnVudHUtcy0xdmNwdS0xZ2ItZnJhMS0wMSIsInN1YiI6InNwbHVua19hZG1pbiIsImF1ZCI6ImdnIiwiaWRwIjoiU3BsdW5rIiwianRpIjoiMzIzOGFhNDA4MDkxNTI5MDkzMDZhNzYxYTk5MWQ5YjEzZjZkNmE5YmI1ZmMzZGM0NTA5MzhmNjY2MDgyODY1NSIsImlhdCI6MTYxNDUwODYyNywiZXhwIjoxNjE3MTAwNjI3LCJuYnIiOjE2MTQ1MDg2Mjd9.UPOpCe3zsi_dZ3P6GomfRGklVL-ef8DyMXH0DAUPMzp3xAUKp_EFxRbguslCbJ0dU1e6O_DpXzzINSaEKlsWqw",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSplunk(nil, m)
			s.currentUser = tt.fields.User
			fmt.Println(s.ListLogs())
		})
	}
}

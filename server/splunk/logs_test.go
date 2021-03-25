package splunk

import (
	"fmt"
	"testing"

	"github.com/bakurits/mattermost-plugin-splunk/server/store"
	"github.com/bakurits/mattermost-plugin-splunk/server/store/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var authToken = "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJzcGx1bmtfYWRtaW4gZnJvbSB1YnVudHUtcy0xdmNwdS0xZ2ItZnJhMS0wMSIsInN1YiI6InNwbHVua19hZG1pbiIsImF1ZCI6Im1hdHRlcm1vc3QgcGx1Z2luIiwiaWRwIjoiU3BsdW5rIiwianRpIjoiM2JkNWZiNThlYmM5NGY5ODUwM2VkNjY0ZDYyMTAyMmM2OGVlMDM3NWU2MjgzMDc0ODcwNjE5OTNmNDJlMjQzOSIsImlhdCI6MTYxNjY2MTAwMiwiZXhwIjoxNjIyNzA5MDAyLCJuYnIiOjE2MTY2NjEwMDJ9.4KEIkoHOBiTtoXZ6lCM-4huE4grzBA2BcFR_MxtU8Vf5rza4lFCEPKfX0_TzNsUFsRIAypy6yLCUgKdYRQ8T4Q"
var server = "https://splunkapi.opsolutions.dev"

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
					Server: server,
					Token:  authToken,
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
					Server: server,
					Token:  authToken,
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

package splunk

import (
	"testing"

	"github.com/mattermost/mattermost-plugin-splunk/server/store"
	"github.com/mattermost/mattermost-plugin-splunk/server/store/mock"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_splunk_ChangeUser(t *testing.T) {
	t.Skip("GH-48 issue link https://github.com/mattermost/mattermost-plugin-splunk/issues/48")
	ctrl := gomock.NewController(t)
	is := assert.New(t)
	m := mock.NewMockStore(ctrl)
	m.EXPECT().ChangeCurrentUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	m.EXPECT().User(gomock.Any(), gomock.Any(), gomock.Any()).Return(store.SplunkUser{}, errors.New("no user found")).AnyTimes()
	defer ctrl.Finish()

	type args struct {
		server string
		id     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "auth success",
			args: args{
				server: server,
				id:     authToken,
			},
			wantErr: false,
		},
		{
			name: "auth failure",
			args: args{
				server: server,
				id:     "dasdas",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, m)
			err := s.LoginUser("", tt.args.server, tt.args.id)
			if tt.wantErr && err == nil {
				is.Fail("expected error but didn't get any")
			}
			if !tt.wantErr && err != nil {
				is.Failf("", "didn't expect any error but got %v", err)
			}
		})
	}
}

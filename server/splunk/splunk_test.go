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

func Test_splunk_extractUserInfo(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		username     string
		token        string
		errorMessage string
	}{
		{
			name:         "id is empty, return error",
			id:           "",
			username:     "",
			token:        "",
			errorMessage: "Please provide username and token like so: username/token. You can user username only if already authenticated",
		},
		{
			name:         "id has more than 2 parameters, return error",
			id:           "johndoe/token/more",
			username:     "",
			token:        "",
			errorMessage: "Arguments to extract username and/or token must be 2",
		},
		{
			name:         "id has only the username parameter, returns valid username",
			id:           "johndoe",
			username:     "johndoe",
			token:        "",
			errorMessage: "",
		},
		{
			name:         "id has username and token parameters, returns username and token",
			id:           "johndoe/token",
			username:     "johndoe",
			token:        "token",
			errorMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, token, err := extractUserInfo(tt.id)

			if tt.errorMessage != "" {
				assert.Equal(t, tt.errorMessage, err.Error())
				return
			}

			assert.Equal(t, tt.username, username)
			assert.Equal(t, tt.token, token)
		})
	}
}

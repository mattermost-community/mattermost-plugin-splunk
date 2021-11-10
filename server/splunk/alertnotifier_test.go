package splunk

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-splunk/server/store/mock"
)

func Test_alertNotifier_delete(t *testing.T) {
	type fields struct {
		alertsInChannel []string
	}
	type args struct {
		data []struct {
			channelID string
			alertID   string
		}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "gg",
			fields: fields{
				alertsInChannel: []string{"1", "2", "3", "4", "5", "6"},
			},
			args: args{
				data: []struct {
					channelID string
					alertID   string
				}{
					{alertID: "1", channelID: "gg"},
					{alertID: "4", channelID: "gg"},
					{alertID: "2", channelID: "gg"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mock.NewMockStore(ctrl)
			defer ctrl.Finish()
			s := newSplunk(nil, m)
			for _, alert := range tt.fields.alertsInChannel {
				m.EXPECT().CreateAlert(tt.name, alert).Return(nil).Times(1)
				if err := s.addAlert(tt.name, alert); (err != nil) != tt.wantErr {
					t.Errorf("addAlertActionFunc() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			for _, arg := range tt.args.data {
				m.EXPECT().DeleteChannelAlert(tt.name, arg.alertID).Return(errors.New("error in deleting alert")).Times(1)
				if err := s.delete(tt.name, arg.alertID); (err != nil) == tt.wantErr {
					t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

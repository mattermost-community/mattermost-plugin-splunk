package splunk

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/mattermost/mattermost-plugin-splunk/server/store/mock"
)

func Test_alertNotifier_delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock.NewMockStore(ctrl)
	m.EXPECT().GetSubscription(gomock.Any()).Return([]string{}, nil).AnyTimes()
	m.EXPECT().SetSubscription(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	defer ctrl.Finish()

	s := newSplunk(nil, m)
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
					{alertID: "6", channelID: "gg"},
					{alertID: "4", channelID: "gg"},
					{alertID: "2", channelID: "gg"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, arg := range tt.fields.alertsInChannel {
				if err := s.addAlertActionFunc("gg", arg); (err != nil) == tt.wantErr {
					t.Errorf("addAlertActionFunc() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			for _, arg := range tt.args.data {
				if err := s.delete(arg.channelID, arg.alertID); (err != nil) != tt.wantErr {
					t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			t.Log(s.list("gg"))
		})
	}
}

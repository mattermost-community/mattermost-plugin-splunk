package splunk

import (
	"sync"
	"testing"
)

func Test_alertNotifier_delete(t *testing.T) {
	type fields struct {
		receivers       map[string]AlertActionFunc
		alertsInChannel map[string][]string
		lock            sync.Locker
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
				receivers: map[string]AlertActionFunc{
					"1": func(payload AlertActionWHPayload) {},
					"2": func(payload AlertActionWHPayload) {},
					"3": func(payload AlertActionWHPayload) {},
					"4": func(payload AlertActionWHPayload) {},
					"5": func(payload AlertActionWHPayload) {},
				},
				alertsInChannel: map[string][]string{
					"gg": {"1", "2", "3", "4", "5"},
				},
				lock: &sync.Mutex{},
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
			s := &splunk{
				notifier: &alertNotifier{
					receivers:       tt.fields.receivers,
					alertsInChannel: tt.fields.alertsInChannel,
					lock:            tt.fields.lock,
				},
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

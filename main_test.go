package main

import (
	"github.com/jarcoal/httpmock"
	"testing"
)

func Test_joinChannel(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type args struct {
		token       string
		channelName string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockResponse func()
	}{
		{
			name: "normal case",
			args: args{
				token:       "dammy",
				channelName: "dammy_channel",
			},
			wantErr: false,
			mockResponse: func() {
				httpmock.RegisterResponder("POST", "https://slack.com/api/conversations.join",
					httpmock.NewStringResponder(200, `{"ok": true}`))
			},
		},
		{
			name: "error case",
			args: args{
				token:       "dammy2",
				channelName: "dammy_channel2",
			},
			wantErr: true,
			mockResponse: func() {
				httpmock.RegisterResponder("POST", "https://slack.com/api/conversations.join",
					httpmock.NewStringResponder(200, `{"ok": false, "error": "channel_not_found"}`))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockResponse()
			if err := joinChannel(tt.args.token, tt.args.channelName); (err != nil) != tt.wantErr {
				t.Errorf("joinChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

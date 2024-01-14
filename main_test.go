package main

import (
	"reflect"
	"testing"
)

func Test_getChannelNamesFromFile(t *testing.T) {
	type args struct {
		filePath string
		want     []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "Valid file",
			args:    args{filePath: "config_test.yaml"},
			want:    []string{"a", "b", "c"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChannelNamesFromFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChannelNamesFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getChannelNamesFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

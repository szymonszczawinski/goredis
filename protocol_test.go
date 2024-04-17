package main

import (
	"reflect"
	"testing"
)

func Test_ParseCommand(t *testing.T) {
	type args struct {
		rawMessage string
	}
	tests := []struct {
		name    string
		args    args
		want    Command
		wantErr bool
	}{
		{"case 1", args{"*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"}, SetCommand{key: []byte("foo"), value: []byte("bar")}, false},
		{"case 2", args{"*4\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommand(tt.args.rawMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

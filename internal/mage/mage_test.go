package mage

import (
	"strings"
	"testing"
)

func TestRunCmdWithResult(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "If valid command given, return result of the command",
			args: args{
				cmd: "ls",
			},
			want:    "mage.go\nmage_test.go\n",
			wantErr: false,
		},
		{
			name: "If empty command given, return result of the command",
			args: args{
				cmd: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunCmdWithResult(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCmdWithResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RunCmdWithResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunCmdWithLog(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "If valid command given, return nil",
			args: args{
				cmd: "ls",
			},
			wantErr: false,
		},
		{
			name: "If command result is empty, shows that it was empty",
			args: args{
				cmd: "echo -n",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunCmdWithLog(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("RunCmdWithLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO
func TestRunLongRunningCmd(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := RunLongRunningCmd(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunLongRunningCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RunLongRunningCmd() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RunLongRunningCmd() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGenerateImageTag(t *testing.T) {
	t.Run("If can geenerate image tag, if git is installed", func(t *testing.T) {
		got, err := GenerateImageTag()
		if err != nil {
			t.Errorf("GenerateImageTag() error = %v, wantErr %v", err, false)
			return
		}
		spl := strings.Split(got, "_")
		if len(spl) < 3 {
			t.Errorf("GenerateImageTag() got = %v, want {branch}_{commit short hash}_{current datetime in UTC}", got)
			return
		}
	})
}

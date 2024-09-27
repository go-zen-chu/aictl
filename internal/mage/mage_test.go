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
				cmd: "echo helloworld",
			},
			want:    "helloworld\n",
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

func TestRunLongRunningCmd(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		wantStderr string
		wantErr    bool
	}{
		{
			name: "If command that takes long time given, return result of the command",
			args: args{
				cmd: `bash -c 'for i in {1..3}; do echo out; echo err >&2; sleep 1; done'`,
			},
			wantStdout: "out\nout\nout\n",
			wantStderr: "err\nerr\nerr\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := RunLongRunningCmdWithLog(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunLongRunningCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantStdout {
				t.Errorf("RunLongRunningCmd() got = %v, want %v", got, tt.wantStdout)
			}
			if got1 != tt.wantStderr {
				t.Errorf("RunLongRunningCmd() got1 = %v, want %v", got1, tt.wantStderr)
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

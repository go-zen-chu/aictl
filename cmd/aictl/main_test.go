package main

import (
	"os"
	"strings"
	"testing"

	"github.com/go-zen-chu/aictl/cmd/aictl/cmd"
)

func TestHelloWorld(t *testing.T) {
	envs := os.Environ()
	args := os.Args
	defer func() {
		for _, env := range envs {
			kv := strings.SplitN(env, "=", 2)
			if err := os.Setenv(kv[0], kv[1]); err != nil {
				t.Errorf("revert original env: %v", err)
			}
		}
		os.Args = args
	}()

	tests := []struct {
		name    string
		envs    map[string]string
		args    []string
		wantErr error
	}{
		{
			name:    "test hello world",
			args:    []string{"aictl", "query", "hello"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("set env: %v", err)
				}
			}
			os.Args = tt.args
			if got := cmd.RootCmdExecute(); got != tt.wantErr {
				t.Errorf("HelloWorld() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

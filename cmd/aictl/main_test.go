package main

import (
	"os"
	"strings"
	"testing"

	"github.com/go-zen-chu/aictl/cmd/aictl/cmd"
)

func TestHelloWorld(t *testing.T) {
	originalEnvs := os.Environ()
	originalArgs := os.Args
	defer func() {
		for _, env := range originalEnvs {
			kv := strings.SplitN(env, "=", 2)
			if err := os.Setenv(kv[0], kv[1]); err != nil {
				t.Errorf("revert original env: %v", err)
			}
		}
		os.Args = originalArgs
	}()

	tests := []struct {
		name       string
		envs       map[string]string
		args       []string
		wantErrMsg string
	}{
		{
			name:       "test hello world",
			args:       []string{"aictl", "query", "hello"},
			wantErrMsg: "AICTL_OPENAI_API_KEY is not set",
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
			if err := cmd.RootCmdExecute(); err.Error() != tt.wantErrMsg {
				t.Errorf("got error message = %v, want %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-zen-chu/aictl/cmd/aictl/cmd"
	"github.com/go-zen-chu/aictl/infra/openai"
	"github.com/go-zen-chu/aictl/internal/di"
	"github.com/go-zen-chu/aictl/usecase/query"
	goa "github.com/sashabaranov/go-openai"
	"go.uber.org/mock/gomock"
)

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "If no args given, it should show help message",
			args: []string{"aictl"},
		},
		{
			name: "If verbose args given, it should show debug level logs",
			args: []string{"aictl", "-v", "help"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &app{
				rootCmd: cmd.NewRootCmd(di.NewContainer()),
			}
			err := app.Run(tt.args)
			if err != nil {
				t.Errorf("got error = %v, want nil", err)
			}
		})
	}
}

func genCommandRequirementsMockWithOpenAIResponse(c *gomock.Controller, responseContent string) cmd.CommandRequirements {
	mcr := cmd.NewMockCommandRequirements(c)
	mgoaic := openai.NewMockGoOpenAIClient(c)
	oaic := openai.NewOpenAIClient(mgoaic)
	uq := query.NewUsecaseQuery(oaic)
	mcr.EXPECT().UsecaseQuery().Return(uq)
	mgoaic.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(
		goa.ChatCompletionResponse{
			Choices: []goa.ChatCompletionChoice{
				{
					Message: goa.ChatCompletionMessage{
						Content: responseContent,
					},
				},
			},
		},
		nil)
	return mcr
}

func TestQueryCmd(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		envs          map[string]string
		beforeEnvVars map[string]string
		mockCmdReq    func(c *gomock.Controller) cmd.CommandRequirements
		wantErr       error
		wantEnvVars   map[string]string
	}{
		{
			name:    "If no query given, it should get an error",
			args:    []string{"aictl", "query"},
			wantErr: errors.New("root command: validation in query: query command requires only 1 argument `query text`"),
		},
		{
			name: "If query given, get response from AI",
			args: []string{"aictl", "query", "hello"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return genCommandRequirementsMockWithOpenAIResponse(c, "Hello! How can I assist you today?\n")
			},
			wantErr: nil,
		},
		// Output format
		{
			name: "If output format set as json, get response from AI with json",
			args: []string{"aictl", "query", "-o", "json", "hello"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return genCommandRequirementsMockWithOpenAIResponse(c, `{
  "message": "Hello! How can I assist you today?"
}
`)
			},
			wantErr: nil,
		},
		{
			name:    "If output format is empty, it should get validation error",
			args:    []string{"aictl", "query", "-o", "", "hello"},
			wantErr: errors.New("root command: validation in query: output format is required but got empty"),
		},
		{
			name:    "If output format set as yaml, it should get error because OpenAI response only supports text or json",
			args:    []string{"aictl", "query", "-o", "yaml", "hello"},
			wantErr: errors.New("root command: validation in query: output format must be text or json but got: yaml"),
		},
		// Language
		{
			name: "If response language given, it should get response from OpenAI with the language",
			args: []string{"aictl", "query", "-l", "Japanese", "hello"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return genCommandRequirementsMockWithOpenAIResponse(c, "こんにちは！いかがですか？何かお手伝いできることがありますか？\n")
			},
			wantErr: nil,
		},
		// Text files
		{
			name: "If valid text files given, get response from OpenAI",
			args: []string{"aictl", "query", "-t", "../../testdata/go_error_sample1.go,../../testdata/go_error_sample2.go", "Could you review this code?"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return genCommandRequirementsMockWithOpenAIResponse(c, "Sure! This is the result...\n")
			},
			wantErr: nil,
		},
		{
			name: "If valid text file with no extension given, get response from OpenAI",
			args: []string{"aictl", "query", "-t", "../../testdata/InvalidDockerfile", "Could you review this code?"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return genCommandRequirementsMockWithOpenAIResponse(c, "Sure! This is the result...\n")
			},
			wantErr: nil,
		},
		{
			name: "If file path is invalid, get an error",
			args: []string{"aictl", "query", "-t", "no_such_file", "Could you review this code?"},
			envs: map[string]string{
				"AICTL_OPENAI_API_KEY": "test",
			},
			wantErr: errors.New("root command: query to openai: filepath no_such_file does not exists: stat no_such_file: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mock
			c := gomock.NewController(t)
			var cr cmd.CommandRequirements
			if tt.mockCmdReq != nil {
				cr = tt.mockCmdReq(c)
			} else {
				cr = di.NewContainer()
			}
			if tt.envs != nil {
				for k, v := range tt.envs {
					if err := os.Setenv(k, v); err != nil {
						t.Errorf("set env var: %s", err)
					}
				}
			}
			// run
			app := &app{
				rootCmd: cmd.NewRootCmd(cr),
			}
			if gotErr := app.Run(tt.args); gotErr != tt.wantErr {
				if gotErr == nil || tt.wantErr == nil {
					t.Errorf("got error = %v, want %v", gotErr, tt.wantErr)
					return
				}
				if gotErr.Error() != tt.wantErr.Error() {
					t.Errorf("got error = %v, want %v", gotErr, tt.wantErr)
				}
			}
		})
	}
}

func TestQueryCmdForGitHubActions(t *testing.T) {
	// create temp dir for test
	tmpDir, err := os.MkdirTemp("", "aictl-*")
	if err != nil {
		t.Errorf("mkdir temp: %s", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("remove all %s: %s", tmpDir, err)
		}
	}()

	// create file for GITHUB_OUTPUT
	ghOutputFilePath := filepath.Join(tmpDir, "githuboutput.txt")
	f, err := os.Create(ghOutputFilePath)
	if err != nil {
		t.Errorf("create file: %s", err)
	}
	defer f.Close()
	if _, err := f.WriteString("test=this is test message\n"); err != nil {
		t.Errorf("write string to file: %s", err)
	}

	// set env vars
	if err := os.Setenv("GITHUB_ACTIONS", "true"); err != nil {
		t.Errorf("set GITHUB_ACTIONS env var: %s", err)
	}
	if err := os.Setenv("GITHUB_OUTPUT", ghOutputFilePath); err != nil {
		t.Errorf("set GITHUB_OUTPUT env var: %s", err)
	}

	tests := []struct {
		name                    string
		args                    []string
		wantErr                 error
		wantGitHubOutputContent string
	}{
		{
			name:                    "If running in GitHub Actions, write result to GITHUB_OUTPUT environment variable",
			args:                    []string{"aictl", "query", "hello"},
			wantErr:                 nil,
			wantGitHubOutputContent: "test=this is test message\nresponse<<AICTL_EOF\nHello! How can I assist you today?\nAICTL_EOF\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			cr := genCommandRequirementsMockWithOpenAIResponse(c, "Hello! How can I assist you today?")

			// run
			app := &app{
				rootCmd: cmd.NewRootCmd(cr),
			}
			if gotErr := app.Run(tt.args); gotErr != tt.wantErr {
				if gotErr == nil || tt.wantErr == nil {
					t.Errorf("got error = %v, want %v", gotErr, tt.wantErr)
					return
				}
				if gotErr.Error() != tt.wantErr.Error() {
					t.Errorf("got error = %v, want %v", gotErr, tt.wantErr)
				}
			}

			// check GITHUB_OUTPUT file
			contentBytes, err := os.ReadFile(ghOutputFilePath)
			if err != nil {
				t.Errorf("read github output file: %s", err)
			}
			if string(contentBytes) != tt.wantGitHubOutputContent {
				t.Errorf("got github output content = %s, want %s", string(contentBytes), tt.wantGitHubOutputContent)
			}
		})
	}
}

package main

import (
	"errors"
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

func TestQueryCmd(t *testing.T) {
	mockCmdReqWithOpenAIResponse := func(c *gomock.Controller, responseContent string) cmd.CommandRequirements {
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

	tests := []struct {
		name       string
		args       []string
		mockCmdReq func(c *gomock.Controller) cmd.CommandRequirements
		wantErr    error
	}{
		{
			name:    "If no query given, it should get an error",
			args:    []string{"aictl", "query"},
			wantErr: errors.New("root command: validation in query: query command requires only 1 argument `query text`"),
		},
		{
			name: "If query given, it should get response from OpenAI",
			args: []string{"aictl", "query", "hello"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return mockCmdReqWithOpenAIResponse(c, "Hello! How can I assist you today?\n")
			},
			wantErr: nil,
		},
		{
			name: "If output format set as json, it should get response from OpenAI with json",
			args: []string{"aictl", "query", "-o", "json", "hello"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return mockCmdReqWithOpenAIResponse(c, `{
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
		{
			name: "If response language given, it should get response from OpenAI with the language",
			args: []string{"aictl", "query", "-l", "Japanese", "hello"},
			mockCmdReq: func(c *gomock.Controller) cmd.CommandRequirements {
				return mockCmdReqWithOpenAIResponse(c, "こんにちは！いかがですか？何かお手伝いできることがありますか？\n")
			},
			wantErr: nil,
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

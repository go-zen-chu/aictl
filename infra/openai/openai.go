package openai

import (
	"context"
	"fmt"

	"github.com/go-zen-chu/aictl/usecase/query"

	goa "github.com/sashabaranov/go-openai"
)

type openaiClient struct {
	cli *goa.Client
}

func NewOpenAIClient(token string) query.OpenAIClient {
	return &openaiClient{
		cli: goa.NewClient(token),
	}
}

var queryTemplate = `%s
Return your response **ONLY** with valid %s format.`

func (c *openaiClient) Ask(ctx context.Context, query string, outputFormat string) (string, error) {
	resType := goa.ChatCompletionResponseFormatTypeText
	if outputFormat == "json" {
		resType = goa.ChatCompletionResponseFormatTypeJSONObject
	}
	resp, err := c.cli.CreateChatCompletion(ctx, goa.ChatCompletionRequest{
		Model: goa.GPT4oMini,
		Messages: []goa.ChatCompletionMessage{
			{
				Role:    goa.ChatMessageRoleUser,
				Content: fmt.Sprintf(queryTemplate, query, outputFormat),
			},
		},
		ResponseFormat: &goa.ChatCompletionResponseFormat{
			Type: resType,
		},
	})
	if err != nil {
		return "", fmt.Errorf("chat completion request: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response from chat completion")
	}
	return resp.Choices[0].Message.Content, nil
}

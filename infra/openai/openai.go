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

func (c *openaiClient) Ask(ctx context.Context, query string) (string, error) {
	resp, err := c.cli.CreateChatCompletion(ctx, goa.ChatCompletionRequest{
		Model: goa.GPT4oMini,
		Messages: []goa.ChatCompletionMessage{
			{
				Role:    goa.ChatMessageRoleUser,
				Content: query,
			},
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

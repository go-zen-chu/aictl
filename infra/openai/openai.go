package openai

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"text/template"

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

func (c *openaiClient) Ask(
	ctx context.Context,
	query string,
	outputFormat string,
	responseLanguage string,
	textFiles []query.InputTextFile) (string, error) {
	resType := goa.ChatCompletionResponseFormatTypeText
	if outputFormat == "json" {
		resType = goa.ChatCompletionResponseFormatTypeJSONObject
	}
	qs := newQueryStruct(query, outputFormat, responseLanguage, textFiles)
	q, err := qs.generateQuery()
	if err != nil {
		return "", fmt.Errorf("generate query: %w", err)
	}
	slog.Debug("Query to OpenAI: %s", q)
	// Ask to OpenAI
	resp, err := c.cli.CreateChatCompletion(ctx, goa.ChatCompletionRequest{
		Model: goa.GPT4oMini,
		Messages: []goa.ChatCompletionMessage{
			{
				Role:    goa.ChatMessageRoleUser,
				Content: q,
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

type queryStruct struct {
	query                string
	outputFormat         string
	responseLanguage     string
	textFiles            []query.InputTextFile
	codeBlockPlaceholder string
}

func newQueryStruct(query string, outputFormat string, responseLanguage string, textFiles []query.InputTextFile) *queryStruct {
	return &queryStruct{
		query:                query,
		outputFormat:         outputFormat,
		responseLanguage:     responseLanguage,
		textFiles:            textFiles,
		codeBlockPlaceholder: "```",
	}
}

var queryTemplate = `{{.query}}

{{range .TextFiles -}}
{{.codeBlockPlaceholder}}{{.Extension}}
{{.Content}}
{{.codeBlockPlaceholder}}
{{end -}}

The following order must be followed:
* Return your response with valid {{.outputFormat}} format only.
* Return your response with {{.outputFormat}} language.`

func (q *queryStruct) generateQuery() (string, error) {
	tmpl, err := template.New("query").Parse(queryTemplate)
	if err != nil {
		return "", fmt.Errorf("parse query template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, q); err != nil {
		return "", fmt.Errorf("execute query template: %w", err)
	}
	return buf.String(), nil
}

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
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

type GoOpenAIClient interface {
	CreateChatCompletion(ctx context.Context, req goa.ChatCompletionRequest) (goa.ChatCompletionResponse, error)
}

type openaiClient struct {
	cli GoOpenAIClient
}

func NewOpenAIClient(cli GoOpenAIClient) query.OpenAIClient {
	return &openaiClient{
		cli: cli,
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
	slog.Debug("Query to OpenAI:", "query", q)
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
	Query            string
	OutputFormat     string
	ResponseLanguage string
	TextFiles        []query.InputTextFile
}

func newQueryStruct(query string, outputFormat string, responseLanguage string, textFiles []query.InputTextFile) *queryStruct {
	return &queryStruct{
		Query:            query,
		OutputFormat:     outputFormat,
		ResponseLanguage: responseLanguage,
		TextFiles:        textFiles,
	}
}

var queryTemplate = `{{.Query}}

{{range .TextFiles -}}
` + "```" + `{{.Extension}}
{{.Content}}
` + "```" + `
{{end -}}

The following order must be followed:
* Return your response with valid {{.OutputFormat}} format only.
* Return your response with {{.ResponseLanguage}} language.`

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

package query

import (
	"context"
	"fmt"
)

type UsecaseQuery interface {
	QueryToOpenAI(query string) (string, error)
}

type OpenAIClient interface {
	Ask(ctx context.Context, query string) (string, error)
}

type usecaseQuery struct {
	openAIClient OpenAIClient
}

func NewUsecaseQuery(openAIClient OpenAIClient) UsecaseQuery {
	return &usecaseQuery{
		openAIClient: openAIClient,
	}
}

func (uq *usecaseQuery) QueryToOpenAI(query string) (string, error) {
	res, err := uq.openAIClient.Ask(context.Background(), query)
	if err != nil {
		return "", fmt.Errorf("ask to openai: %w", err)
	}
	return res, nil
}

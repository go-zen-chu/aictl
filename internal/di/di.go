package di

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-zen-chu/aictl/infra/openai"
	"github.com/go-zen-chu/aictl/usecase/query"
)

type Container struct {
	cache map[string]any
}

func NewContainer() *Container {
	return &Container{
		cache: map[string]any{},
	}
}

func initOnce[T any](c *Container, component string, fn func() (T, error)) T {
	if v, ok := c.cache[component]; ok {
		return v.(T)
	}
	var err error
	v, err := fn()
	if err != nil {
		slog.Error("failed to set up "+component, "error", err)
		os.Exit(1)
	}
	c.cache[component] = v
	return v
}

func (c *Container) UsecaseQuery() query.UsecaseQuery {
	return initOnce(c, "UsecaseQuery", func() (query.UsecaseQuery, error) {
		return query.NewUsecaseQuery(c.OpenAIClient()), nil
	})
}

func (c *Container) OpenAIClient() query.OpenAIClient {
	return initOnce(c, "OpenAIClient", func() (query.OpenAIClient, error) {
		k := os.Getenv("AICTL_OPENAI_API_KEY")
		if k == "" {
			return nil, fmt.Errorf("AICTL_OPENAI_API_KEY is not set")
		}
		return openai.NewOpenAIClient(k), nil
	})
}

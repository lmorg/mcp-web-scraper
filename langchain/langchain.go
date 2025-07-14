package langchain

import (
	"context"
	"time"

	"github.com/lmorg/mcp-scrape-page/internal"
	"github.com/tmc/langchaingo/callbacks"
)

type Scraper struct {
	CallbacksHandler callbacks.Handler
}

func New() *Scraper {
	return &Scraper{}
}

func (*Scraper) Name() string {
	return internal.Name
}

func (*Scraper) Description() string {
	return internal.Description
}

func (*Scraper) Call(ctx context.Context, input string) (string, error) {
	timeout, _ := context.WithTimeout(ctx, 10*time.Second)
	return internal.Scrape(timeout, input)
}

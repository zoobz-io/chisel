package chisel

import (
	"context"
	"fmt"
)

// Chunker routes content to the appropriate language provider.
type Chunker struct {
	providers map[Language]Provider
}

// New creates a Chunker with the given providers.
func New(providers ...Provider) *Chunker {
	c := &Chunker{
		providers: make(map[Language]Provider),
	}
	for _, p := range providers {
		c.providers[p.Language()] = p
	}
	return c
}

// Register adds a provider to the chunker.
func (c *Chunker) Register(p Provider) {
	c.providers[p.Language()] = p
}

// Chunk parses content using the appropriate provider for the language.
func (c *Chunker) Chunk(ctx context.Context, lang Language, filename string, content []byte) ([]Chunk, error) {
	p, ok := c.providers[lang]
	if !ok {
		return nil, fmt.Errorf("no provider for language: %s", lang)
	}
	return p.Chunk(ctx, filename, content)
}

// Languages returns all registered languages.
func (c *Chunker) Languages() []Language {
	langs := make([]Language, 0, len(c.providers))
	for lang := range c.providers {
		langs = append(langs, lang)
	}
	return langs
}

// HasProvider returns true if a provider is registered for the language.
func (c *Chunker) HasProvider(lang Language) bool {
	_, ok := c.providers[lang]
	return ok
}

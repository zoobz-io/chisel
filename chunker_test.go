package chisel

import (
	"context"
	"testing"
)

// mockProvider implements Provider for testing.
type mockProvider struct {
	lang   Language
	chunks []Chunk
	err    error
}

func (m *mockProvider) Language() Language {
	return m.lang
}

func (m *mockProvider) Chunk(_ context.Context, _ string, _ []byte) ([]Chunk, error) {
	return m.chunks, m.err
}

func TestNew(t *testing.T) {
	p1 := &mockProvider{lang: Go}
	p2 := &mockProvider{lang: Python}

	c := New(p1, p2)

	if !c.HasProvider(Go) {
		t.Error("expected Go provider to be registered")
	}
	if !c.HasProvider(Python) {
		t.Error("expected Python provider to be registered")
	}
	if c.HasProvider(Rust) {
		t.Error("expected Rust provider to not be registered")
	}
}

func TestChunker_Register(t *testing.T) {
	c := New()

	if c.HasProvider(Go) {
		t.Error("expected no Go provider initially")
	}

	c.Register(&mockProvider{lang: Go})

	if !c.HasProvider(Go) {
		t.Error("expected Go provider after registration")
	}
}

func TestChunker_Chunk(t *testing.T) {
	expected := []Chunk{
		{Symbol: "TestFunc", Kind: KindFunction, StartLine: 1, EndLine: 5},
	}

	c := New(&mockProvider{
		lang:   Go,
		chunks: expected,
	})

	chunks, err := c.Chunk(context.Background(), Go, "test.go", []byte("package main"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) != len(expected) {
		t.Errorf("got %d chunks, want %d", len(chunks), len(expected))
	}
}

func TestChunker_Chunk_NoProvider(t *testing.T) {
	c := New()

	_, err := c.Chunk(context.Background(), Go, "test.go", []byte("package main"))
	if err == nil {
		t.Error("expected error for missing provider")
	}
}

func TestChunker_Languages(t *testing.T) {
	c := New(
		&mockProvider{lang: Go},
		&mockProvider{lang: Python},
		&mockProvider{lang: Rust},
	)

	langs := c.Languages()
	if len(langs) != 3 {
		t.Errorf("got %d languages, want 3", len(langs))
	}
}

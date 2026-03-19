package markdown

import (
	"context"
	"testing"

	"github.com/zoobz-io/chisel"
)

func TestProvider_Language(t *testing.T) {
	p := New()
	if p.Language() != chisel.Markdown {
		t.Errorf("Language() = %v, want %v", p.Language(), chisel.Markdown)
	}
}

func TestProvider_Chunk_SingleSection(t *testing.T) {
	src := `# Getting Started

This is the introduction.

Some more text here.
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "README.md", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) != 1 {
		t.Fatalf("got %d chunks, want 1", len(chunks))
	}

	chunk := chunks[0]
	if chunk.Symbol != "Getting Started" {
		t.Errorf("Symbol = %q, want %q", chunk.Symbol, "Getting Started")
	}
	if chunk.Kind != chisel.KindSection {
		t.Errorf("Kind = %v, want %v", chunk.Kind, chisel.KindSection)
	}
}

func TestProvider_Chunk_MultipleSections(t *testing.T) {
	src := `# Title

Intro text.

## Installation

Install instructions.

## Usage

Usage instructions.
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "README.md", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) != 3 {
		t.Fatalf("got %d chunks, want 3", len(chunks))
	}

	symbols := []string{"Title", "Installation", "Usage"}
	for i, want := range symbols {
		if chunks[i].Symbol != want {
			t.Errorf("chunks[%d].Symbol = %q, want %q", i, chunks[i].Symbol, want)
		}
	}
}

func TestProvider_Chunk_NestedContext(t *testing.T) {
	src := `# API

## Methods

### Get

Get method description.

### Post

Post method description.
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "api.md", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) < 4 {
		t.Fatalf("got %d chunks, want at least 4", len(chunks))
	}

	// Find "Get" section and check context
	var getChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Symbol == "Get" {
			getChunk = &chunks[i]
			break
		}
	}

	if getChunk == nil {
		t.Fatal("Get chunk not found")
	}

	// Context should include parent headers
	if len(getChunk.Context) < 1 {
		t.Errorf("expected context, got %v", getChunk.Context)
	}
}

func TestProvider_Chunk_EmptyDocument(t *testing.T) {
	p := New()
	chunks, err := p.Chunk(context.Background(), "empty.md", []byte(""))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) != 0 {
		t.Errorf("got %d chunks, want 0", len(chunks))
	}
}

func TestProvider_Chunk_NoHeaders(t *testing.T) {
	src := `Just some text without any headers.

More text here.
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "notes.md", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should have one chunk with empty symbol
	if len(chunks) != 1 {
		t.Fatalf("got %d chunks, want 1", len(chunks))
	}
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		input     string
		wantLevel int
		wantTitle string
	}{
		{"# Title", 1, "Title"},
		{"## Subtitle", 2, "Subtitle"},
		{"### Deep", 3, "Deep"},
		{"###### Level 6", 6, "Level 6"},
		{"####### Too Deep", 0, ""}, // Invalid: > 6
		{"Not a header", 0, ""},
		{"#NoSpace", 1, "NoSpace"},
	}

	for _, tt := range tests {
		level, title := parseHeader(tt.input)
		if level != tt.wantLevel {
			t.Errorf("parseHeader(%q) level = %d, want %d", tt.input, level, tt.wantLevel)
		}
		if title != tt.wantTitle {
			t.Errorf("parseHeader(%q) title = %q, want %q", tt.input, title, tt.wantTitle)
		}
	}
}

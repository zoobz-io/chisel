package golang

import (
	"context"
	"testing"

	"github.com/zoobzio/chisel"
)

func TestProvider_Language(t *testing.T) {
	p := New()
	if p.Language() != chisel.Go {
		t.Errorf("Language() = %v, want %v", p.Language(), chisel.Go)
	}
}

func TestProvider_Chunk_Function(t *testing.T) {
	src := `package main

// Add adds two numbers.
func Add(a, b int) int {
	return a + b
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.go", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) != 1 {
		t.Fatalf("got %d chunks, want 1", len(chunks))
	}

	chunk := chunks[0]
	if chunk.Symbol != "Add" {
		t.Errorf("Symbol = %q, want %q", chunk.Symbol, "Add")
	}
	if chunk.Kind != chisel.KindFunction {
		t.Errorf("Kind = %v, want %v", chunk.Kind, chisel.KindFunction)
	}
}

func TestProvider_Chunk_Method(t *testing.T) {
	src := `package main

type Calculator struct{}

// Add adds two numbers.
func (c *Calculator) Add(a, b int) int {
	return a + b
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.go", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should have: type Calculator, method Add
	if len(chunks) < 2 {
		t.Fatalf("got %d chunks, want at least 2", len(chunks))
	}

	// Find the method chunk
	var methodChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Kind == chisel.KindMethod {
			methodChunk = &chunks[i]
			break
		}
	}

	if methodChunk == nil {
		t.Fatal("no method chunk found")
	}

	if methodChunk.Symbol != "Calculator.Add" {
		t.Errorf("Symbol = %q, want %q", methodChunk.Symbol, "Calculator.Add")
	}
	if len(methodChunk.Context) == 0 {
		t.Error("expected Context to contain receiver type")
	}
}

func TestProvider_Chunk_Type(t *testing.T) {
	src := `package main

// User represents a user.
type User struct {
	Name string
	Age  int
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.go", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) != 1 {
		t.Fatalf("got %d chunks, want 1", len(chunks))
	}

	chunk := chunks[0]
	if chunk.Symbol != "User" {
		t.Errorf("Symbol = %q, want %q", chunk.Symbol, "User")
	}
	if chunk.Kind != chisel.KindClass {
		t.Errorf("Kind = %v, want %v", chunk.Kind, chisel.KindClass)
	}
}

func TestProvider_Chunk_Interface(t *testing.T) {
	src := `package main

// Reader reads data.
type Reader interface {
	Read(p []byte) (n int, err error)
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.go", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) != 1 {
		t.Fatalf("got %d chunks, want 1", len(chunks))
	}

	chunk := chunks[0]
	if chunk.Symbol != "Reader" {
		t.Errorf("Symbol = %q, want %q", chunk.Symbol, "Reader")
	}
	if chunk.Kind != chisel.KindInterface {
		t.Errorf("Kind = %v, want %v", chunk.Kind, chisel.KindInterface)
	}
}

func TestProvider_Chunk_ParseError(t *testing.T) {
	src := `package main

func broken( {
`
	p := New()
	_, err := p.Chunk(context.Background(), "test.go", []byte(src))
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestProvider_Chunk_PackageDoc(t *testing.T) {
	src := `// Package main is the entry point.
package main

func main() {}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.go", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should have package doc and main function
	if len(chunks) < 2 {
		t.Fatalf("got %d chunks, want at least 2", len(chunks))
	}

	// First chunk should be package doc
	if chunks[0].Kind != chisel.KindModule {
		t.Errorf("first chunk Kind = %v, want %v", chunks[0].Kind, chisel.KindModule)
	}
}

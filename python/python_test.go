package python

import (
	"context"
	"testing"

	"github.com/zoobz-io/chisel"
)

func TestProvider_Language(t *testing.T) {
	p := New()
	if p.Language() != chisel.Python {
		t.Errorf("Language() = %v, want %v", p.Language(), chisel.Python)
	}
}

func TestProvider_Chunk_Function(t *testing.T) {
	src := `def add(a, b):
    """Add two numbers."""
    return a + b
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.py", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	var funcChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Kind == chisel.KindFunction {
			funcChunk = &chunks[i]
			break
		}
	}

	if funcChunk == nil {
		t.Fatal("no function chunk found")
	}

	if funcChunk.Symbol != "add" {
		t.Errorf("Symbol = %q, want %q", funcChunk.Symbol, "add")
	}
}

func TestProvider_Chunk_Class(t *testing.T) {
	src := `class Calculator:
    """A simple calculator."""

    def add(self, a, b):
        return a + b

    def subtract(self, a, b):
        return a - b
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "calc.py", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should have: class, add method, subtract method
	if len(chunks) < 3 {
		t.Fatalf("got %d chunks, want at least 3", len(chunks))
	}

	// Find class chunk
	var classChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Kind == chisel.KindClass {
			classChunk = &chunks[i]
			break
		}
	}

	if classChunk == nil {
		t.Fatal("no class chunk found")
	}

	if classChunk.Symbol != "Calculator" {
		t.Errorf("Symbol = %q, want %q", classChunk.Symbol, "Calculator")
	}

	// Find method chunks and verify context
	var methodCount int
	for i := range chunks {
		if chunks[i].Kind == chisel.KindMethod {
			methodCount++
			if len(chunks[i].Context) == 0 {
				t.Errorf("method %q has no context", chunks[i].Symbol)
			}
		}
	}

	if methodCount < 2 {
		t.Errorf("got %d methods, want at least 2", methodCount)
	}
}

func TestProvider_Chunk_NestedClass(t *testing.T) {
	src := `class Outer:
    class Inner:
        def method(self):
            pass
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "nested.py", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should have nested classes
	var classCount int
	for i := range chunks {
		if chunks[i].Kind == chisel.KindClass {
			classCount++
		}
	}

	if classCount < 2 {
		t.Errorf("got %d classes, want at least 2", classCount)
	}
}

func TestProvider_Chunk_Decorator(t *testing.T) {
	src := `@decorator
def decorated_function():
    pass
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "decorated.py", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	// Decorated function should still be extracted
	var hasFunc bool
	for i := range chunks {
		if chunks[i].Symbol == "decorated_function" {
			hasFunc = true
			break
		}
	}

	if !hasFunc {
		t.Error("decorated function not found")
	}
}

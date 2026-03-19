package rust

import (
	"context"
	"testing"

	"github.com/zoobz-io/chisel"
)

func TestProvider_Language(t *testing.T) {
	p := New()
	if p.Language() != chisel.Rust {
		t.Errorf("Language() = %v, want %v", p.Language(), chisel.Rust)
	}
}

func TestProvider_Chunk_Function(t *testing.T) {
	src := `fn add(a: i32, b: i32) -> i32 {
    a + b
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "lib.rs", []byte(src))
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

func TestProvider_Chunk_Struct(t *testing.T) {
	src := `struct User {
    name: String,
    age: u32,
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "lib.rs", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	var structChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Kind == chisel.KindType && chunks[i].Symbol == "User" {
			structChunk = &chunks[i]
			break
		}
	}

	if structChunk == nil {
		t.Fatal("no struct chunk found")
	}
}

func TestProvider_Chunk_Impl(t *testing.T) {
	src := `struct Calculator;

impl Calculator {
    fn add(&self, a: i32, b: i32) -> i32 {
        a + b
    }

    fn subtract(&self, a: i32, b: i32) -> i32 {
        a - b
    }
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "calc.rs", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should have: struct, impl block, and methods
	if len(chunks) < 3 {
		t.Fatalf("got %d chunks, want at least 3", len(chunks))
	}

	// Find methods and verify context
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

func TestProvider_Chunk_Trait(t *testing.T) {
	src := `trait Drawable {
    fn draw(&self);
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "lib.rs", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	var traitChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Symbol == "Drawable" {
			traitChunk = &chunks[i]
			break
		}
	}

	if traitChunk == nil {
		t.Fatal("no trait chunk found")
	}

	if traitChunk.Kind != chisel.KindInterface {
		t.Errorf("Kind = %v, want %v", traitChunk.Kind, chisel.KindInterface)
	}
}

func TestProvider_Chunk_Enum(t *testing.T) {
	src := `enum Status {
    Active,
    Inactive,
    Pending,
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "lib.rs", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	var enumChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Symbol == "Status" {
			enumChunk = &chunks[i]
			break
		}
	}

	if enumChunk == nil {
		t.Fatal("no enum chunk found")
	}

	if enumChunk.Kind != chisel.KindEnum {
		t.Errorf("Kind = %v, want %v", enumChunk.Kind, chisel.KindEnum)
	}
}

func TestProvider_Chunk_Module(t *testing.T) {
	src := `mod utils {
    fn helper() {}
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "lib.rs", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	var modChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Symbol == "utils" && chunks[i].Kind == chisel.KindModule {
			modChunk = &chunks[i]
			break
		}
	}

	if modChunk == nil {
		t.Fatal("no module chunk found")
	}
}

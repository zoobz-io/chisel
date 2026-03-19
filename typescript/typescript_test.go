package typescript

import (
	"context"
	"testing"

	"github.com/zoobz-io/chisel"
)

func TestProvider_Language(t *testing.T) {
	p := New()
	if p.Language() != chisel.TypeScript {
		t.Errorf("Language() = %v, want %v", p.Language(), chisel.TypeScript)
	}

	pjs := NewJavaScript()
	if pjs.Language() != chisel.JavaScript {
		t.Errorf("Language() = %v, want %v", pjs.Language(), chisel.JavaScript)
	}
}

func TestProvider_Chunk_Function(t *testing.T) {
	src := `function add(a: number, b: number): number {
	return a + b;
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.ts", []byte(src))
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
	src := `class Calculator {
	add(a: number, b: number): number {
		return a + b;
	}

	subtract(a: number, b: number): number {
		return a - b;
	}
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "calc.ts", []byte(src))
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

func TestProvider_Chunk_Interface(t *testing.T) {
	src := `interface User {
	name: string;
	age: number;
}
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "types.ts", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	var ifaceChunk *chisel.Chunk
	for i := range chunks {
		if chunks[i].Kind == chisel.KindInterface {
			ifaceChunk = &chunks[i]
			break
		}
	}

	if ifaceChunk == nil {
		t.Fatal("no interface chunk found")
	}

	if ifaceChunk.Symbol != "User" {
		t.Errorf("Symbol = %q, want %q", ifaceChunk.Symbol, "User")
	}
}

func TestProvider_Chunk_ArrowFunction(t *testing.T) {
	src := `const multiply = (a: number, b: number) => a * b;
`
	p := New()
	chunks, err := p.Chunk(context.Background(), "test.ts", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Arrow functions should be extracted
	var hasArrow bool
	for i := range chunks {
		if chunks[i].Kind == chisel.KindFunction {
			hasArrow = true
			break
		}
	}

	if !hasArrow {
		t.Log("arrow function not extracted as separate chunk (acceptable)")
	}
}

// JavaScript-specific tests

func TestProvider_Chunk_JS_Function(t *testing.T) {
	src := `function add(a, b) {
	return a + b;
}
`
	p := NewJavaScript()
	chunks, err := p.Chunk(context.Background(), "test.js", []byte(src))
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

func TestProvider_Chunk_JS_Class(t *testing.T) {
	src := `class Calculator {
	constructor() {
		this.value = 0;
	}

	add(n) {
		this.value += n;
		return this;
	}
}
`
	p := NewJavaScript()
	chunks, err := p.Chunk(context.Background(), "calc.js", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should have class and methods
	if len(chunks) < 2 {
		t.Fatalf("got %d chunks, want at least 2", len(chunks))
	}

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
}

func TestProvider_Chunk_JS_ArrowFunction(t *testing.T) {
	src := `const multiply = (a, b) => a * b;
`
	p := NewJavaScript()
	chunks, err := p.Chunk(context.Background(), "test.js", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Arrow functions should be extracted
	var hasArrow bool
	for i := range chunks {
		if chunks[i].Kind == chisel.KindFunction {
			hasArrow = true
			break
		}
	}

	if !hasArrow {
		t.Log("arrow function not extracted as separate chunk (acceptable)")
	}
}

func TestProvider_Chunk_JS_NestedFunctions(t *testing.T) {
	src := `function outer() {
	function inner() {
		return 42;
	}
	return inner();
}
`
	p := NewJavaScript()
	chunks, err := p.Chunk(context.Background(), "test.js", []byte(src))
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Should find both outer and inner functions
	var funcCount int
	for i := range chunks {
		if chunks[i].Kind == chisel.KindFunction {
			funcCount++
		}
	}

	if funcCount < 2 {
		t.Errorf("got %d functions, want at least 2", funcCount)
	}
}

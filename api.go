// Package chisel provides AST-aware code chunking for semantic search and embeddings.
package chisel

import "context"

// Language identifies a programming language.
type Language string

// Supported languages.
const (
	Go         Language = "go"
	TypeScript Language = "typescript"
	JavaScript Language = "javascript"
	Python     Language = "python"
	Rust       Language = "rust"
	Markdown   Language = "markdown"
)

// Kind categorizes a chunk.
type Kind string

// Chunk kinds.
const (
	KindFunction  Kind = "function"
	KindMethod    Kind = "method"
	KindClass     Kind = "class"
	KindInterface Kind = "interface"
	KindType      Kind = "type"
	KindEnum      Kind = "enum"
	KindConstant  Kind = "constant"
	KindVariable  Kind = "variable"
	KindSection   Kind = "section" // For markdown headers
	KindModule    Kind = "module"  // Package/file level
)

// Chunk represents a semantic unit of code or documentation.
type Chunk struct {
	// Content is the actual code or text.
	Content string

	// Symbol is the name of the function, class, type, or section.
	Symbol string

	// Kind categorizes this chunk.
	Kind Kind

	// StartLine is the 1-indexed starting line number.
	StartLine int

	// EndLine is the 1-indexed ending line number.
	EndLine int

	// Context is the parent chain for this chunk.
	// Example: ["class UserService", "method getUser"]
	Context []string
}

// Provider parses a specific language into chunks.
type Provider interface {
	// Chunk parses content and returns semantic chunks.
	Chunk(ctx context.Context, filename string, content []byte) ([]Chunk, error)

	// Language returns the supported language.
	Language() Language
}

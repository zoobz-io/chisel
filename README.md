# chisel

[![CI Status](https://github.com/zoobz-io/chisel/workflows/CI/badge.svg)](https://github.com/zoobz-io/chisel/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zoobz-io/chisel/graph/badge.svg?branch=main)](https://codecov.io/gh/zoobz-io/chisel)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoobz-io/chisel)](https://goreportcard.com/report/github.com/zoobz-io/chisel)
[![CodeQL](https://github.com/zoobz-io/chisel/workflows/CodeQL/badge.svg)](https://github.com/zoobz-io/chisel/security/code-scanning)
[![Go Reference](https://pkg.go.dev/badge/github.com/zoobz-io/chisel.svg)](https://pkg.go.dev/github.com/zoobz-io/chisel)
[![License](https://img.shields.io/github/license/zoobz-io/chisel)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zoobz-io/chisel)](go.mod)
[![Release](https://img.shields.io/github/v/release/zoobz-io/chisel)](https://github.com/zoobz-io/chisel/releases)

AST-aware code chunking for semantic search and embeddings. Chisel parses source code into meaningful units—functions, classes, methods—preserving the context that makes code searchable.

## From Syntax to Semantics

```go
source := []byte(`
func New(cfg Config) *Handler { ... }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { ... }

type Config struct {
    Timeout time.Duration
    Logger  *slog.Logger
}
`)

chunks, _ := c.Chunk(ctx, chisel.Go, "api.go", source)

for _, chunk := range chunks {
    fmt.Printf("[%s] %s (lines %d-%d)\n", chunk.Kind, chunk.Symbol, chunk.StartLine, chunk.EndLine)
}
// [function] New (lines 2-2)
// [method] Handler.ServeHTTP (lines 4-4)
// [class] Config (lines 6-9)
```

Every chunk carries its symbol name, kind, line range, and parent context. Methods know their receiver. Nested types know their enclosing scope.

```go
chunk := chunks[1]
// chunk.Symbol    → "Handler.ServeHTTP"
// chunk.Kind      → "method"
// chunk.Context   → ["Handler"]
// chunk.Content   → the full method source
// chunk.StartLine → 4
// chunk.EndLine   → 4
```

Feed chunks to an embedding model, store in a vector database, and search code by meaning rather than text.

## Install

```bash
go get github.com/zoobz-io/chisel
```

**Language providers** (install only what you need):

```bash
go get github.com/zoobz-io/chisel/golang     # Go (stdlib, no deps)
go get github.com/zoobz-io/chisel/markdown   # Markdown (no deps)
go get github.com/zoobz-io/chisel/typescript # TypeScript/JavaScript (tree-sitter)
go get github.com/zoobz-io/chisel/python     # Python (tree-sitter)
go get github.com/zoobz-io/chisel/rust       # Rust (tree-sitter)
```

Requires Go 1.24+.

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/zoobz-io/chisel"
    "github.com/zoobz-io/chisel/golang"
    "github.com/zoobz-io/chisel/typescript"
)

func main() {
    // Create a chunker with language providers
    c := chisel.New(
        golang.New(),
        typescript.New(),
        typescript.NewJavaScript(),
    )

    source := []byte(`
package auth

// Authenticate validates user credentials.
func Authenticate(username, password string) (*User, error) {
    // ...
}

// User represents an authenticated user.
type User struct {
    ID    string
    Email string
}
`)

    chunks, err := c.Chunk(context.Background(), chisel.Go, "auth.go", source)
    if err != nil {
        panic(err)
    }

    for _, chunk := range chunks {
        fmt.Printf("[%s] %s\n", chunk.Kind, chunk.Symbol)
        fmt.Printf("  Lines: %d-%d\n", chunk.StartLine, chunk.EndLine)
        if len(chunk.Context) > 0 {
            fmt.Printf("  Context: %v\n", chunk.Context)
        }
    }
}
```

Output:

```text
[function] Authenticate
  Lines: 4-6
[class] User
  Lines: 8-12
```

## Capabilities

| Feature | Description | Docs |
|---------|-------------|------|
| **Multi-language** | Go, TypeScript, JavaScript, Python, Rust, Markdown | [Providers](docs/2.guides/1.providers.md) |
| **Semantic extraction** | Functions, methods, classes, interfaces, types, enums | [Concepts](docs/1.learn/3.concepts.md) |
| **Context preservation** | Parent chain for nested definitions | [Architecture](docs/1.learn/4.architecture.md) |
| **Line mapping** | Precise source locations for each chunk | [Types](docs/4.reference/2.types.md) |
| **Zero-copy providers** | Go and Markdown use stdlib only | [Architecture](docs/1.learn/4.architecture.md) |

## Why Chisel?

- **Semantic boundaries** — Chunks split at function/class boundaries, not arbitrary line counts
- **Embedding-ready** — Output designed for vector databases and semantic search
- **Isolated dependencies** — Tree-sitter only where needed; Go/Markdown have zero external deps
- **Context-aware** — Methods know their parent class; nested functions know their scope
- **Consistent interface** — Same `Provider` contract across all languages

## Code Intelligence Pipelines

Chisel enables a pattern: **parse once, search by meaning**.

Your codebase becomes a corpus of semantic units. Each function, method, and type gets embedded with its full context — symbol name, parent scope, documentation. Queries match intent, not just text.

```go
// Chunk your codebase
chunks, _ := c.Chunk(ctx, chisel.Go, path, source)

// Embed each chunk (using your embedding provider)
for _, chunk := range chunks {
    embedding := embedder.Embed(chunk.Content)
    vectorDB.Store(embedding, chunk.Symbol, chunk.Kind, path)
}

// Search by meaning
results := vectorDB.Query("authentication middleware")
// Returns: AuthMiddleware, ValidateToken, SessionHandler
// Not just files containing the word "authentication"
```

Symbol names and kinds become metadata. Line ranges enable source navigation. Context chains power hierarchical search.

## Ecosystem

Chisel provides the chunking layer for code intelligence pipelines:

- **[vicky](https://github.com/zoobz-io/vicky)** — Code search and retrieval service

## Documentation

- **Learn**
  - [Overview](docs/1.learn/1.overview.md) — What chisel is and why
  - [Quickstart](docs/1.learn/2.quickstart.md) — Get productive in minutes
  - [Concepts](docs/1.learn/3.concepts.md) — Core abstractions
  - [Architecture](docs/1.learn/4.architecture.md) — How it works internally
- **Guides**
  - [Providers](docs/2.guides/1.providers.md) — Language-specific details
  - [Testing](docs/2.guides/2.testing.md) — Testing code that uses chisel
  - [Troubleshooting](docs/2.guides/3.troubleshooting.md) — Common issues
- **Reference**
  - [API](docs/4.reference/1.api.md) — Function signatures
  - [Types](docs/4.reference/2.types.md) — Type definitions

## Contributing

Contributions welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT — see [LICENSE](LICENSE) for details.

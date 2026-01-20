# Benchmarks

Performance benchmarks for chisel providers.

## Running

```bash
make test-bench

# Or directly:
go test -bench=. github.com/zoobzio/chisel/testing/benchmarks -benchmem
```

## Current Results

Representative results on AMD Ryzen 5 3600X (~50-line files):

| Provider | Time | Memory | Allocations |
|----------|------|--------|-------------|
| Go | 32µs | 17KB | 402 |
| TypeScript | 313µs | 63KB | 579 |
| Python | 328µs | 63KB | 569 |
| Rust | 293µs | 61KB | 566 |
| Markdown | 4µs | 7KB | 45 |

## Analysis

- **Go provider** uses stdlib `go/parser`, ~10x faster than tree-sitter
- **Markdown** is fastest (simple string scanning, no AST)
- **Tree-sitter providers** have similar performance characteristics

## Adding Benchmarks

Follow the pattern in `benchmarks_test.go`:

```go
func BenchmarkNewProvider(b *testing.B) {
    p := newprovider.New()
    ctx := context.Background()
    source := []byte(`...`)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := p.Chunk(ctx, "file.ext", source)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

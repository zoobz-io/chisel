# Testing

Test utilities and infrastructure for chisel.

## Structure

```text
testing/
├── helpers.go          # Test assertion helpers
├── helpers_test.go     # Tests for helpers
├── benchmarks/         # Performance benchmarks
└── integration/        # End-to-end tests
```

## Helpers

The `testing` package provides domain-specific assertion helpers:

```go
import chitesting "github.com/zoobzio/chisel/testing"

func TestMyChunker(t *testing.T) {
    chunks := getChunks()

    chitesting.AssertChunkCount(t, chunks, 3)
    chitesting.AssertHasSymbol(t, chunks, "MyFunction")
    chitesting.AssertHasKind(t, chunks, chisel.KindFunction)
}
```

### Available Helpers

| Function | Description |
|----------|-------------|
| `AssertChunkCount(t, chunks, n)` | Assert exact chunk count |
| `AssertHasSymbol(t, chunks, sym)` | Assert symbol exists |
| `AssertHasKind(t, chunks, kind)` | Assert kind exists |
| `FindBySymbol(chunks, sym)` | Find chunk by symbol |
| `FindByKind(chunks, kind)` | Find chunk by kind |
| `CountByKind(chunks, kind)` | Count chunks by kind |

## Running Tests

```bash
# All tests
make test

# Unit tests only (short mode)
make test-unit

# Integration tests
make test-integration

# Benchmarks
make test-bench
```

## Coverage

```bash
make coverage
```

Target: 70% project, 80% patch.

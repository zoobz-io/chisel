# Integration Tests

End-to-end tests for chisel.

## Purpose

Integration tests verify that:

- Providers work correctly with real-world code samples
- The `Chunker` correctly routes to providers
- Multi-file scenarios produce expected results

## Running

```bash
make test-integration

# Or directly:
go test -v ./testing/integration/...
```

## Structure

```text
integration/
├── README.md
├── testdata/           # Sample source files
│   ├── go/
│   ├── typescript/
│   ├── python/
│   └── rust/
└── integration_test.go
```

## Adding Tests

1. Add sample files to `testdata/[language]/`
2. Add test cases in `integration_test.go`
3. Verify expected chunks match actual output

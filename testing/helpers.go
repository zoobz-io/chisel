// Package testing provides test utilities for chisel.
package testing

import (
	"testing"

	"github.com/zoobzio/chisel"
)

// AssertChunkCount checks that the expected number of chunks were produced.
func AssertChunkCount(t *testing.T, chunks []chisel.Chunk, want int) {
	t.Helper()
	if len(chunks) != want {
		t.Errorf("got %d chunks, want %d", len(chunks), want)
	}
}

// AssertHasSymbol checks that a chunk with the given symbol exists.
func AssertHasSymbol(t *testing.T, chunks []chisel.Chunk, symbol string) {
	t.Helper()
	for _, c := range chunks {
		if c.Symbol == symbol {
			return
		}
	}
	t.Errorf("no chunk with symbol %q found", symbol)
}

// AssertHasKind checks that at least one chunk has the given kind.
func AssertHasKind(t *testing.T, chunks []chisel.Chunk, kind chisel.Kind) {
	t.Helper()
	for _, c := range chunks {
		if c.Kind == kind {
			return
		}
	}
	t.Errorf("no chunk with kind %v found", kind)
}

// FindBySymbol returns the first chunk with the given symbol.
func FindBySymbol(chunks []chisel.Chunk, symbol string) *chisel.Chunk {
	for i := range chunks {
		if chunks[i].Symbol == symbol {
			return &chunks[i]
		}
	}
	return nil
}

// FindByKind returns the first chunk with the given kind.
func FindByKind(chunks []chisel.Chunk, kind chisel.Kind) *chisel.Chunk {
	for i := range chunks {
		if chunks[i].Kind == kind {
			return &chunks[i]
		}
	}
	return nil
}

// CountByKind returns the number of chunks with the given kind.
func CountByKind(chunks []chisel.Chunk, kind chisel.Kind) int {
	count := 0
	for _, c := range chunks {
		if c.Kind == kind {
			count++
		}
	}
	return count
}

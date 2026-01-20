package testing

import (
	"testing"

	"github.com/zoobzio/chisel"
)

func TestFindBySymbol(t *testing.T) {
	chunks := []chisel.Chunk{
		{Symbol: "foo", Kind: chisel.KindFunction},
		{Symbol: "bar", Kind: chisel.KindClass},
	}

	result := FindBySymbol(chunks, "foo")
	if result == nil {
		t.Error("expected to find foo")
	}
	if result.Symbol != "foo" {
		t.Errorf("Symbol = %q, want %q", result.Symbol, "foo")
	}

	result = FindBySymbol(chunks, "baz")
	if result != nil {
		t.Error("expected not to find baz")
	}
}

func TestFindByKind(t *testing.T) {
	chunks := []chisel.Chunk{
		{Symbol: "foo", Kind: chisel.KindFunction},
		{Symbol: "bar", Kind: chisel.KindClass},
	}

	result := FindByKind(chunks, chisel.KindClass)
	if result == nil {
		t.Error("expected to find class")
	}
	if result.Kind != chisel.KindClass {
		t.Errorf("Kind = %v, want %v", result.Kind, chisel.KindClass)
	}
}

func TestCountByKind(t *testing.T) {
	chunks := []chisel.Chunk{
		{Kind: chisel.KindFunction},
		{Kind: chisel.KindFunction},
		{Kind: chisel.KindClass},
	}

	count := CountByKind(chunks, chisel.KindFunction)
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}

	count = CountByKind(chunks, chisel.KindMethod)
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
}

func TestAssertChunkCount(t *testing.T) {
	chunks := []chisel.Chunk{
		{Symbol: "foo"},
		{Symbol: "bar"},
	}

	// Success case - should not fail
	AssertChunkCount(t, chunks, 2)
}

func TestAssertHasSymbol(t *testing.T) {
	chunks := []chisel.Chunk{
		{Symbol: "foo", Kind: chisel.KindFunction},
		{Symbol: "bar", Kind: chisel.KindClass},
	}

	// Success case - should not fail
	AssertHasSymbol(t, chunks, "foo")
	AssertHasSymbol(t, chunks, "bar")
}

func TestAssertHasKind(t *testing.T) {
	chunks := []chisel.Chunk{
		{Symbol: "foo", Kind: chisel.KindFunction},
		{Symbol: "bar", Kind: chisel.KindClass},
	}

	// Success case - should not fail
	AssertHasKind(t, chunks, chisel.KindFunction)
	AssertHasKind(t, chunks, chisel.KindClass)
}

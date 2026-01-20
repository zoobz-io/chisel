package benchmarks

import (
	"context"
	"testing"

	"github.com/zoobzio/chisel/golang"
	"github.com/zoobzio/chisel/markdown"
	"github.com/zoobzio/chisel/python"
	"github.com/zoobzio/chisel/rust"
	"github.com/zoobzio/chisel/typescript"
)

// Sample sources for benchmarking

var goSource = []byte(`package main

import "fmt"

// Calculator performs arithmetic operations.
type Calculator struct {
	value int
}

// New creates a new Calculator.
func New() *Calculator {
	return &Calculator{}
}

// Add adds n to the current value.
func (c *Calculator) Add(n int) *Calculator {
	c.value += n
	return c
}

// Subtract subtracts n from the current value.
func (c *Calculator) Subtract(n int) *Calculator {
	c.value -= n
	return c
}

// Multiply multiplies the current value by n.
func (c *Calculator) Multiply(n int) *Calculator {
	c.value *= n
	return c
}

// Value returns the current value.
func (c *Calculator) Value() int {
	return c.value
}

func main() {
	calc := New()
	result := calc.Add(10).Multiply(2).Subtract(5).Value()
	fmt.Println(result)
}
`)

var tsSource = []byte(`interface Calculator {
	add(n: number): Calculator;
	subtract(n: number): Calculator;
	value(): number;
}

class BasicCalculator implements Calculator {
	private _value: number;

	constructor() {
		this._value = 0;
	}

	add(n: number): Calculator {
		this._value += n;
		return this;
	}

	subtract(n: number): Calculator {
		this._value -= n;
		return this;
	}

	multiply(n: number): Calculator {
		this._value *= n;
		return this;
	}

	value(): number {
		return this._value;
	}
}

function createCalculator(): Calculator {
	return new BasicCalculator();
}

const calc = createCalculator();
const result = calc.add(10).subtract(5).value();
console.log(result);
`)

var pySource = []byte(`"""Calculator module for arithmetic operations."""

class Calculator:
    """A simple calculator class."""

    def __init__(self):
        """Initialize the calculator with value 0."""
        self._value = 0

    def add(self, n: int) -> "Calculator":
        """Add n to the current value."""
        self._value += n
        return self

    def subtract(self, n: int) -> "Calculator":
        """Subtract n from the current value."""
        self._value -= n
        return self

    def multiply(self, n: int) -> "Calculator":
        """Multiply the current value by n."""
        self._value *= n
        return self

    @property
    def value(self) -> int:
        """Return the current value."""
        return self._value


def create_calculator() -> Calculator:
    """Create a new calculator instance."""
    return Calculator()


if __name__ == "__main__":
    calc = create_calculator()
    result = calc.add(10).subtract(5).value
    print(result)
`)

var rustSource = []byte(`//! Calculator module for arithmetic operations.

/// A simple calculator struct.
pub struct Calculator {
    value: i32,
}

impl Calculator {
    /// Creates a new Calculator with value 0.
    pub fn new() -> Self {
        Calculator { value: 0 }
    }

    /// Adds n to the current value.
    pub fn add(&mut self, n: i32) -> &mut Self {
        self.value += n;
        self
    }

    /// Subtracts n from the current value.
    pub fn subtract(&mut self, n: i32) -> &mut Self {
        self.value -= n;
        self
    }

    /// Multiplies the current value by n.
    pub fn multiply(&mut self, n: i32) -> &mut Self {
        self.value *= n;
        self
    }

    /// Returns the current value.
    pub fn value(&self) -> i32 {
        self.value
    }
}

fn main() {
    let mut calc = Calculator::new();
    calc.add(10).subtract(5);
    println!("{}", calc.value());
}
`)

var mdSource = []byte(`# Calculator Documentation

This document describes the Calculator API.

## Overview

The Calculator provides basic arithmetic operations.

## API Reference

### Constructor

Creates a new calculator instance with initial value 0.

### Methods

#### add(n)

Adds n to the current value and returns the calculator for chaining.

#### subtract(n)

Subtracts n from the current value and returns the calculator for chaining.

#### multiply(n)

Multiplies the current value by n and returns the calculator for chaining.

#### value()

Returns the current calculated value.

## Examples

Here's a simple example:

` + "```" + `go
calc := New()
result := calc.Add(10).Multiply(2).Value()
` + "```" + `

## Conclusion

The Calculator is a simple but powerful tool for chaining arithmetic operations.
`)

func BenchmarkGolang(b *testing.B) {
	p := golang.New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Chunk(ctx, "calc.go", goSource)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTypeScript(b *testing.B) {
	p := typescript.New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Chunk(ctx, "calc.ts", tsSource)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPython(b *testing.B) {
	p := python.New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Chunk(ctx, "calc.py", pySource)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRust(b *testing.B) {
	p := rust.New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Chunk(ctx, "calc.rs", rustSource)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdown(b *testing.B) {
	p := markdown.New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Chunk(ctx, "calc.md", mdSource)
		if err != nil {
			b.Fatal(err)
		}
	}
}

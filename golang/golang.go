// Package golang provides Go source code chunking using the standard library parser.
package golang

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/zoobz-io/chisel"
)

// Provider chunks Go source files using go/parser.
type Provider struct{}

// New creates a new Go chunking provider.
func New() *Provider {
	return &Provider{}
}

// Language returns the Go language identifier.
func (p *Provider) Language() chisel.Language {
	return chisel.Go
}

// Chunk parses Go source and extracts semantic chunks.
func (p *Provider) Chunk(_ context.Context, filename string, content []byte) ([]chisel.Chunk, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	var chunks []chisel.Chunk

	// Extract package-level documentation
	if file.Doc != nil {
		chunks = append(chunks, chisel.Chunk{
			Content:   file.Doc.Text(),
			Symbol:    file.Name.Name,
			Kind:      chisel.KindModule,
			StartLine: fset.Position(file.Doc.Pos()).Line,
			EndLine:   fset.Position(file.Doc.End()).Line,
		})
	}

	// Extract declarations
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			chunk := extractFunc(fset, d, content)
			chunks = append(chunks, chunk)

		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						chunk := extractType(fset, d, ts, content)
						chunks = append(chunks, chunk)
					}
				}
			}
		}
	}

	return chunks, nil
}

// extractFunc extracts a function or method declaration.
func extractFunc(fset *token.FileSet, fn *ast.FuncDecl, content []byte) chisel.Chunk {
	start := fset.Position(fn.Pos())
	end := fset.Position(fn.End())

	var doc string
	if fn.Doc != nil {
		doc = fn.Doc.Text()
	}

	// Determine symbol name and kind
	name := fn.Name.Name
	kind := chisel.KindFunction
	var ctx []string

	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		kind = chisel.KindMethod
		recvType := receiverType(fn.Recv.List[0].Type)
		if recvType != "" {
			ctx = append(ctx, "type "+recvType)
			name = recvType + "." + name
		}
	}

	// Extract the full function source
	src := safeSlice(content, int(fn.Pos()-1), int(fn.End()-1))

	return chisel.Chunk{
		Content:   doc + string(src),
		Symbol:    name,
		Kind:      kind,
		StartLine: start.Line,
		EndLine:   end.Line,
		Context:   ctx,
	}
}

// extractType extracts a type declaration.
func extractType(fset *token.FileSet, decl *ast.GenDecl, ts *ast.TypeSpec, content []byte) chisel.Chunk {
	start := fset.Position(decl.Pos())
	end := fset.Position(decl.End())

	var doc string
	if decl.Doc != nil {
		doc = decl.Doc.Text()
	}

	// Determine kind based on type
	kind := chisel.KindType
	if _, ok := ts.Type.(*ast.StructType); ok {
		kind = chisel.KindClass
	} else if _, ok := ts.Type.(*ast.InterfaceType); ok {
		kind = chisel.KindInterface
	}

	src := safeSlice(content, int(decl.Pos()-1), int(decl.End()-1))

	return chisel.Chunk{
		Content:   doc + string(src),
		Symbol:    ts.Name.Name,
		Kind:      kind,
		StartLine: start.Line,
		EndLine:   end.Line,
	}
}

// receiverType extracts the type name from a receiver expression.
func receiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.Ident:
		return t.Name
	}
	return ""
}

// safeSlice safely extracts a slice from content.
func safeSlice(content []byte, start, end int) []byte {
	if start < 0 {
		start = 0
	}
	if end > len(content) {
		end = len(content)
	}
	if start >= end {
		return nil
	}
	return content[start:end]
}

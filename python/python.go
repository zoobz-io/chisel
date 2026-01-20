// Package python provides Python chunking using tree-sitter.
package python

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"

	"github.com/zoobzio/chisel"
)

// Provider chunks Python files using tree-sitter.
type Provider struct{}

// New creates a new Python chunking provider.
func New() *Provider {
	return &Provider{}
}

// Language returns the Python language identifier.
func (p *Provider) Language() chisel.Language {
	return chisel.Python
}

// Chunk parses Python source and extracts semantic chunks.
func (p *Provider) Chunk(_ context.Context, filename string, content []byte) ([]chisel.Chunk, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(python.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	var chunks []chisel.Chunk
	walkNode(tree.RootNode(), content, nil, &chunks)

	return chunks, nil
}

// walkNode recursively walks the AST and extracts chunks.
func walkNode(node *sitter.Node, content []byte, ctx []string, chunks *[]chisel.Chunk) {
	nodeType := node.Type()

	switch nodeType {
	case "function_definition":
		chunk := extractFunction(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "class_definition":
		chunk := extractClass(node, content, ctx)
		*chunks = append(*chunks, chunk)

		// Walk children with class context
		className := getChildByField(node, "name", content)
		newCtx := append(ctx, "class "+className)
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "block" {
				walkNode(child, content, newCtx, chunks)
			}
		}
		return
	}

	// Walk children
	for i := 0; i < int(node.ChildCount()); i++ {
		walkNode(node.Child(i), content, ctx, chunks)
	}
}

// extractFunction extracts a function definition.
func extractFunction(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)

	// Determine if this is a method (has class context)
	kind := chisel.KindFunction
	if len(ctx) > 0 {
		kind = chisel.KindMethod
	}

	// Extract docstring if present
	var docstring string
	body := node.ChildByFieldName("body")
	if body != nil && body.ChildCount() > 0 {
		first := body.Child(0)
		if first.Type() == "expression_statement" && first.ChildCount() > 0 {
			expr := first.Child(0)
			if expr.Type() == "string" {
				docstring = string(content[expr.StartByte():expr.EndByte()])
			}
		}
	}

	fullContent := string(content[node.StartByte():node.EndByte()])
	if docstring != "" {
		fullContent = docstring + "\n" + fullContent
	}

	return chisel.Chunk{
		Content:   fullContent,
		Symbol:    name,
		Kind:      kind,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// extractClass extracts a class definition.
func extractClass(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      chisel.KindClass,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// getChildByField finds a named child and returns its content.
func getChildByField(node *sitter.Node, field string, content []byte) string {
	child := node.ChildByFieldName(field)
	if child == nil {
		return ""
	}
	return string(content[child.StartByte():child.EndByte()])
}

// copyContext creates a copy of the context slice.
func copyContext(ctx []string) []string {
	if ctx == nil {
		return nil
	}
	result := make([]string, len(ctx))
	copy(result, ctx)
	return result
}

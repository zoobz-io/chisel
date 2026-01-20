// Package typescript provides TypeScript/JavaScript chunking using tree-sitter.
package typescript

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"

	"github.com/zoobzio/chisel"
)

// Provider chunks TypeScript files using tree-sitter.
type Provider struct {
	lang chisel.Language
}

// New creates a new TypeScript chunking provider.
func New() *Provider {
	return &Provider{lang: chisel.TypeScript}
}

// NewJavaScript creates a provider configured for JavaScript files.
func NewJavaScript() *Provider {
	return &Provider{lang: chisel.JavaScript}
}

// Language returns the language identifier.
func (p *Provider) Language() chisel.Language {
	return p.lang
}

// Chunk parses TypeScript/JavaScript source and extracts semantic chunks.
func (p *Provider) Chunk(_ context.Context, _ string, content []byte) ([]chisel.Chunk, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(typescript.GetLanguage())

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
	case "function_declaration", "arrow_function", "function":
		chunk := extractFunction(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "method_definition":
		chunk := extractMethod(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "class_declaration":
		chunk := extractClass(node, content, ctx)
		*chunks = append(*chunks, chunk)

		// Walk children with class context
		className := getChildByField(node, "name", content)
		newCtx := append(copyContext(ctx), "class "+className)
		for i := 0; i < int(node.ChildCount()); i++ {
			walkNode(node.Child(i), content, newCtx, chunks)
		}
		return

	case "interface_declaration":
		chunk := extractInterface(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "type_alias_declaration":
		chunk := extractType(node, content, ctx)
		*chunks = append(*chunks, chunk)
	}

	// Walk children
	for i := 0; i < int(node.ChildCount()); i++ {
		walkNode(node.Child(i), content, ctx, chunks)
	}
}

// extractFunction extracts a function declaration.
func extractFunction(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)
	if name == "" {
		name = "<anonymous>"
	}

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      chisel.KindFunction,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// extractMethod extracts a method definition.
func extractMethod(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)
	if name == "" {
		name = "<anonymous>"
	}

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      chisel.KindMethod,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// extractClass extracts a class declaration.
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

// extractInterface extracts an interface declaration.
func extractInterface(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      chisel.KindInterface,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// extractType extracts a type alias.
func extractType(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      chisel.KindType,
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

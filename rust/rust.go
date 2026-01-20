// Package rust provides Rust chunking using tree-sitter.
package rust

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"

	"github.com/zoobzio/chisel"
)

// Provider chunks Rust files using tree-sitter.
type Provider struct{}

// New creates a new Rust chunking provider.
func New() *Provider {
	return &Provider{}
}

// Language returns the Rust language identifier.
func (p *Provider) Language() chisel.Language {
	return chisel.Rust
}

// Chunk parses Rust source and extracts semantic chunks.
func (p *Provider) Chunk(_ context.Context, _ string, content []byte) ([]chisel.Chunk, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(rust.GetLanguage())

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
	case "function_item":
		chunk := extractFunction(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "impl_item":
		// Extract impl block and walk its methods
		chunk := extractImpl(node, content, ctx)
		*chunks = append(*chunks, chunk)

		// Get the type being implemented
		typeName := getImplTypeName(node, content)
		newCtx := copyContext(ctx)
		if typeName != "" {
			newCtx = append(newCtx, "impl "+typeName)
		}

		// Walk children with impl context
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "declaration_list" {
				walkNode(child, content, newCtx, chunks)
			}
		}
		return

	case "struct_item":
		chunk := extractStruct(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "enum_item":
		chunk := extractEnum(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "trait_item":
		chunk := extractTrait(node, content, ctx)
		*chunks = append(*chunks, chunk)

	case "mod_item":
		chunk := extractModule(node, content, ctx)
		*chunks = append(*chunks, chunk)
	}

	// Walk children
	for i := 0; i < int(node.ChildCount()); i++ {
		walkNode(node.Child(i), content, ctx, chunks)
	}
}

// extractFunction extracts a function item.
func extractFunction(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)

	kind := chisel.KindFunction
	if len(ctx) > 0 {
		kind = chisel.KindMethod
	}

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      kind,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// extractImpl extracts an impl block.
func extractImpl(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	typeName := getImplTypeName(node, content)

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    typeName,
		Kind:      chisel.KindClass,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// extractStruct extracts a struct definition.
func extractStruct(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
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

// extractEnum extracts an enum definition.
func extractEnum(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      chisel.KindEnum,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// extractTrait extracts a trait definition.
func extractTrait(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
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

// extractModule extracts a module definition.
func extractModule(node *sitter.Node, content []byte, ctx []string) chisel.Chunk {
	name := getChildByField(node, "name", content)

	return chisel.Chunk{
		Content:   string(content[node.StartByte():node.EndByte()]),
		Symbol:    name,
		Kind:      chisel.KindModule,
		StartLine: int(node.StartPoint().Row) + 1,
		EndLine:   int(node.EndPoint().Row) + 1,
		Context:   copyContext(ctx),
	}
}

// getImplTypeName extracts the type name from an impl block.
func getImplTypeName(node *sitter.Node, content []byte) string {
	// Look for type_identifier child
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "type_identifier" {
			return string(content[child.StartByte():child.EndByte()])
		}
		if child.Type() == "generic_type" {
			// Get the base type from generic
			typeIdent := child.ChildByFieldName("type")
			if typeIdent != nil {
				return string(content[typeIdent.StartByte():typeIdent.EndByte()])
			}
		}
	}
	return ""
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

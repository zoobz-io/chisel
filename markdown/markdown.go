// Package markdown provides Markdown document chunking by header sections.
package markdown

import (
	"context"
	"strings"

	"github.com/zoobz-io/chisel"
)

// Provider chunks Markdown files by header sections.
type Provider struct{}

// New creates a new Markdown chunking provider.
func New() *Provider {
	return &Provider{}
}

// Language returns the Markdown language identifier.
func (p *Provider) Language() chisel.Language {
	return chisel.Markdown
}

// Chunk splits Markdown content into sections based on headers.
func (p *Provider) Chunk(_ context.Context, _ string, content []byte) ([]chisel.Chunk, error) {
	text := string(content)
	lines := strings.Split(text, "\n")

	var chunks []chisel.Chunk
	var currentSection strings.Builder
	var sectionStart int
	var sectionName string
	var sectionLevel int
	var context []string

	flush := func(endLine int) {
		content := strings.TrimSpace(currentSection.String())
		if content != "" {
			chunks = append(chunks, chisel.Chunk{
				Content:   content,
				Symbol:    sectionName,
				Kind:      chisel.KindSection,
				StartLine: sectionStart,
				EndLine:   endLine,
				Context:   copyContext(context),
			})
		}
		currentSection.Reset()
	}

	for i, line := range lines {
		lineNum := i + 1

		// Check for ATX-style headers (# Header)
		if strings.HasPrefix(line, "#") {
			level, title := parseHeader(line)
			if level > 0 {
				// Flush previous section
				flush(lineNum - 1)

				// Update context based on header level
				context = updateContext(context, sectionLevel, level, sectionName)

				sectionStart = lineNum
				sectionName = title
				sectionLevel = level
			}
		}

		currentSection.WriteString(line)
		currentSection.WriteString("\n")
	}

	// Flush final section
	flush(len(lines))

	return chunks, nil
}

// parseHeader extracts the level and title from a header line.
func parseHeader(line string) (level int, title string) {
	for _, ch := range line {
		if ch == '#' {
			level++
		} else {
			break
		}
	}

	if level == 0 || level > 6 {
		return 0, ""
	}

	title = strings.TrimSpace(strings.TrimLeft(line, "# "))
	return level, title
}

// updateContext adjusts the context chain when entering a new header.
func updateContext(ctx []string, prevLevel, newLevel int, prevName string) []string {
	if prevLevel == 0 {
		return nil
	}

	// Build context from previous sections
	if newLevel > prevLevel && prevName != "" {
		// Deeper level: add previous as parent
		return append(ctx, prevName)
	} else if newLevel <= prevLevel {
		// Same or shallower: pop context to appropriate level
		depth := newLevel - 1
		if depth < 0 {
			depth = 0
		}
		if depth < len(ctx) {
			return ctx[:depth]
		}
	}

	return ctx
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

package chisel

import "testing"

func TestLanguageConstants(t *testing.T) {
	tests := []struct {
		lang Language
		want string
	}{
		{Go, "go"},
		{TypeScript, "typescript"},
		{JavaScript, "javascript"},
		{Python, "python"},
		{Rust, "rust"},
		{Markdown, "markdown"},
	}

	for _, tt := range tests {
		if string(tt.lang) != tt.want {
			t.Errorf("Language %v = %q, want %q", tt.lang, string(tt.lang), tt.want)
		}
	}
}

func TestKindConstants(t *testing.T) {
	tests := []struct {
		kind Kind
		want string
	}{
		{KindFunction, "function"},
		{KindMethod, "method"},
		{KindClass, "class"},
		{KindType, "type"},
		{KindSection, "section"},
		{KindModule, "module"},
	}

	for _, tt := range tests {
		if string(tt.kind) != tt.want {
			t.Errorf("Kind %v = %q, want %q", tt.kind, string(tt.kind), tt.want)
		}
	}
}

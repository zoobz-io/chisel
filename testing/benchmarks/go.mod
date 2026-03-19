module github.com/zoobz-io/chisel/testing/benchmarks

go 1.24.5

require (
	github.com/zoobz-io/chisel/golang v0.0.0
	github.com/zoobz-io/chisel/markdown v0.0.0
	github.com/zoobz-io/chisel/python v0.0.0
	github.com/zoobz-io/chisel/rust v0.0.0
	github.com/zoobz-io/chisel/typescript v0.0.0
)

require (
	github.com/smacker/go-tree-sitter v0.0.0-20240827094217-dd81d9e9be82 // indirect
	github.com/zoobz-io/chisel v0.0.0 // indirect
)

replace github.com/zoobz-io/chisel => ../../

replace github.com/zoobz-io/chisel/golang => ../../golang

replace github.com/zoobz-io/chisel/markdown => ../../markdown

replace github.com/zoobz-io/chisel/python => ../../python

replace github.com/zoobz-io/chisel/rust => ../../rust

replace github.com/zoobz-io/chisel/typescript => ../../typescript

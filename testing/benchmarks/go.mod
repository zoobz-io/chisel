module github.com/zoobzio/chisel/testing/benchmarks

go 1.24.5

require (
	github.com/zoobzio/chisel/golang v0.0.0
	github.com/zoobzio/chisel/markdown v0.0.0
	github.com/zoobzio/chisel/python v0.0.0
	github.com/zoobzio/chisel/rust v0.0.0
	github.com/zoobzio/chisel/typescript v0.0.0
)

require (
	github.com/smacker/go-tree-sitter v0.0.0-20240827094217-dd81d9e9be82 // indirect
	github.com/zoobzio/chisel v0.0.0 // indirect
)

replace github.com/zoobzio/chisel => ../../

replace github.com/zoobzio/chisel/golang => ../../golang

replace github.com/zoobzio/chisel/markdown => ../../markdown

replace github.com/zoobzio/chisel/python => ../../python

replace github.com/zoobzio/chisel/rust => ../../rust

replace github.com/zoobzio/chisel/typescript => ../../typescript

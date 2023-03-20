
main.wasm: main.go
	tinygo build -no-debug -panic=trap -gc=leaking -target wasm main.go

wasm-opt: main.wasm
	wasm-opt -O main.wasm -o main.wasm --enable-bulk-memory

default: main.wasm


# Prover wasm compilation
compile-prover-go-groth16:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o wasm/circuit.wasm wasm/groth16/main.go
	wasm-opt -O wasm/circuit.wasm -o wasm/circuit.wasm --enable-bulk-memory

compile-prover-go-plonk:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o wasm/circuit.wasm wasm/plonk/main.go
	wasm-opt -O wasm/circuit.wasm -o wasm/circuit.wasm --enable-bulk-memory

compile-prover-tinygo-groth16:
	tinygo build -no-debug -panic=trap -gc=leaking -target=wasm -o wasm/circuit.wasm wasm/groth16/main.go
	wasm-opt -O wasm/circuit.wasm -o wasm/circuit.wasm --enable-bulk-memory

compile-prover-tinygo-plonk:
	tinygo build -no-debug -panic=trap -gc=leaking -target=wasm -o wasm/circuit.wasm wasm/plonk/main.go
	wasm-opt -O wasm/circuit.wasm -o wasm/circuit.wasm --enable-bulk-memory

# Gnark circuit compilation
compile-circuit-groth16:
	@go run ./cmd/compiler -backend groth16

compile-circuit-plonk:
	@go run ./cmd/compiler -backend plonk

# Examples
run-go-example:
	@go run examples/gowasm/main.go

run-tinygo-example:
	@go run examples/tinygo/main.go
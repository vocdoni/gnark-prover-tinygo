# Prover wasm compilation
compile-prover-go-groth16:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o wasm/circuit.wasm wasm/groth16/main.go
	wasm-opt -O wasm/circuit.wasm -o wasm/circuit.wasm --enable-bulk-memory

compile-prover-go-plonk:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o wasm/circuit.wasm wasm/plonk/main.go
	wasm-opt -O wasm/circuit.wasm -o wasm/circuit.wasm --enable-bulk-memory

compile-prover-tinygo-groth16:
	tinygo build -no-debug -panic=trap -gc=leaking -target=wasm -o examples/tinygowasm/circuit.wasm wasm/groth16/main_tinygo.go
	wasm-opt -O examples/tinygowasm/circuit.wasm -o examples/tinygowasm/circuit.wasm --enable-bulk-memory

compile-prover-tinygo-plonk:
	tinygo build -no-debug -panic=trap -gc=leaking -target=wasm -o examples/tinygowasm/circuit.wasm wasm/plonk/main_tinygo.go
	wasm-opt -O examples/tinygowasm/circuit.wasm -o wasm/circuit.wasm --enable-bulk-memory

# Gnark circuit compilation
compile-circuit-groth16:
	@go run ./cmd/compiler -backend groth16

compile-circuit-plonk:
	@go run ./cmd/compiler -backend plonk

# Examples
run-go-example:
	@go run examples/gowasm/main.go

run-tinygo-example:
	@go run examples/tinygowasm/main.go

# MiMC vs. Poseidon test
run-mimc-poseidon-test:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o test/artifacts/groth16-circuit.wasm wasm/groth16/main.go
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o test/artifacts/plonk-circuit.wasm wasm/plonk/main.go

	wasm-opt -O test/artifacts/groth16-circuit.wasm -o test/artifacts/groth16-circuit.wasm --enable-bulk-memory
	wasm-opt -O test/artifacts/plonk-circuit.wasm -o test/artifacts/plonk-circuit.wasm --enable-bulk-memory

	@go test -v ./test
	@go run test/main.go
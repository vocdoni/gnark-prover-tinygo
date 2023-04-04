compile-prover-go:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o artifacts/prover.wasm wasm/main.go
	wasm-opt -O artifacts/prover.wasm -o artifacts/prover.wasm --enable-bulk-memory

compile-prover-tinygo:
	tinygo build -target=wasm -o artifacts/prover.wasm wasm/main.go
#	wasm-opt -O artifacts/circuit.wasm -o artifacts/circuit.wasm --enable-bulk-memory

compile-circuit:
	@go run ./cmd/compiler

run-go-web-example:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit
	@echo "copying artifacts"
	@cp ./artifacts/zkcensus.ccs ./wasm/zkcensus.ccs
	@cp ./artifacts/zkcensus.srs ./wasm/zkcensus.srs
	@cp ./artifacts/zkcensus.pkey ./wasm/zkcensus.pkey
	@echo "compilling the prover for go-wasm"
	@make compile-prover-go
	@echo "removing copied artifacts"
	@rm ./wasm/zkcensus.ccs
	@rm ./wasm/zkcensus.srs
	@rm ./wasm/zkcensus.pkey
	@echo "copying compatible wasm_exec.js"
	@cp ./artifacts/wasm_exec_go.js ./examples/web/wasm_exec.js
	@go run examples/web/main.go

run-tinygo-web-example:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit
	@echo "copying artifacts"
	@cp ./artifacts/zkcensus.ccs ./wasm/zkcensus.ccs
	@cp ./artifacts/zkcensus.srs ./wasm/zkcensus.srs
	@cp ./artifacts/zkcensus.pkey ./wasm/zkcensus.pkey
	@echo "compilling the prover for tinygo"
	@make compile-prover-tinygo
	@echo "removing copied artifacts"
	@rm ./wasm/zkcensus.ccs
	@rm ./wasm/zkcensus.srs
	@rm ./wasm/zkcensus.pkey
	@echo "copying compatible wasm_exec.js"
	@cp ./artifacts/wasm_exec_tinygo.js ./examples/web/wasm_exec.js
	@go run examples/web/main.go

compile-prover-go:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o artifacts/prover.wasm wasm/main.go
	wasm-opt -O artifacts/prover.wasm -o artifacts/prover.wasm --enable-bulk-memory

compile-prover-tinygo-wasi:
	tinygo build -target=wasi -o examples/tinygowasi/artifacts/prover.wasm wasi/main.go

compile-prover-tinygo:
	tinygo build -target=wasm -opt=1 -no-debug -scheduler=asyncify -o artifacts/g16_prover.wasm wasm/main.go
	wasm-opt -O artifacts/g16_prover.wasm -o artifacts/g16_prover.wasm --enable-bulk-memory

compile-prover-g16-tinygo:
	tinygo build -target=wasm -opt=1 -scheduler=asyncify -o artifacts/g16_prover.wasm wasm/g16/main.go
	wasm-opt -O artifacts/g16_prover.wasm -o artifacts/g16_prover.wasm --enable-bulk-memory

compile-circuit:
	@go run ./cmd/compiler --protocol=plonk

compile-circuit-g16:
	@go run ./cmd/compiler --protocol=groth16


run-go-example:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit
	@echo "copying artifacts"
	@go run ./cmd/prover/main.go -circuit=artifacts/zkcensus.ccs -pkey=artifacts/zkcensus.pkey -srs=artifacts/zkcensus.srs -witness=artifacts/zkcensus.witness

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

run-wasi-web-example:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit
	@echo "copying artifacts"
	@cp ./artifacts/zkcensus.ccs ./wasi/zkcensus.ccs
	@cp ./artifacts/zkcensus.srs ./wasi/zkcensus.srs
	@cp ./artifacts/zkcensus.pkey ./wasi/zkcensus.pkey
	@cp ./artifacts/zkcensus.witness ./examples/tinygowasi/artifacts/zkcensus.witness
	@echo "compilling the prover for tinygo (wasi)"
	@make compile-prover-tinygo-wasi
	@echo "removing copied artifacts"
	@rm ./wasi/zkcensus.ccs
	@rm ./wasi/zkcensus.srs
	@rm ./wasi/zkcensus.pkey
	@cd ./examples/tinygowasi && npm i && npx parcel index.html

run-tinygo-web-example-g16:
	@echo "compilling circuit and genering artifacts for groth16"
	@make compile-circuit-g16
	@echo "copying artifacts"
	@cp ./artifacts/g16_zkcensus.ccs ./wasm/g16/zkcensus.ccs
	@cp ./artifacts/g16_zkcensus.pkey ./wasm/g16/zkcensus.pkey
	@echo "compilling the prover for tinygo"
	@make compile-prover-g16-tinygo
	@echo "removing copied artifacts"
	@rm ./wasm/g16/zkcensus.ccs
	@rm ./wasm/g16/zkcensus.pkey
	@echo "copying compatible wasm_exec.js"
	@cp ./artifacts/wasm_exec_tinygo.js ./examples/web/g16/wasm_exec.js
	@go run examples/web/g16/main.go

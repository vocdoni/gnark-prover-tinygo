compile-prover-go-plonk:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o artifacts/plonk_prover.wasm wasm/plonk/main.go
	wasm-opt -O artifacts/plonk_prover.wasm -o artifacts/plonk_prover.wasm --enable-bulk-memory

compile-prover-go-g16:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o artifacts/g16_prover.wasm wasm/g16/main.go
	wasm-opt -O artifacts/g16_prover.wasm -o artifacts/g16_prover.wasm --enable-bulk-memory

compile-prover-tinygo-wasi:
	tinygo build -target=wasi -o examples/tinygowasi/artifacts/prover.wasm wasi/main.go

compile-prover-tinygo-plonk:
	tinygo build -target=wasm -opt=1 -no-debug -scheduler=asyncify -o artifacts/plonk_prover.wasm wasm/plonk/main.go
	wasm-opt -O artifacts/plonk_prover.wasm -o artifacts/plonk_prover.wasm --enable-bulk-memory

compile-prover-tinygo-g16:
	tinygo build -target=wasm -no-debug -opt=1 -scheduler=asyncify -o artifacts/g16_prover.wasm wasm/g16/main.go
	wasm-opt -O artifacts/g16_prover.wasm -o artifacts/g16_prover.wasm --enable-bulk-memory

compile-circuit-plonk:
	@go run ./cmd/compiler --protocol=plonk

compile-circuit-g16:
	@go run ./cmd/compiler --protocol=groth16

run-go-example-plonk:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit-plonk
	@echo "copying artifacts"
	@go run ./cmd/prover/main.go --protocol=plonk -circuit=artifacts/zkcensus.ccs -pkey=artifacts/zkcensus.pkey -srs=artifacts/zkcensus.srs -witness=artifacts/zkcensus.witness

run-go-example-g16:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit-g16
	@echo "copying artifacts"
	@go run ./cmd/prover/main.go --protocol=groth16 -circuit=artifacts/g16_zkcensus.ccs -pkey=artifacts/g16_zkcensus.pkey -witness=artifacts/zkcensus.witness

run-go-web-example-plonk:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit-plonk
	@echo "copying artifacts"
	@cp ./artifacts/zkcensus.ccs ./wasm/plonk/zkcensus.ccs
	@cp ./artifacts/zkcensus.srs ./wasm/plonk/zkcensus.srs
	@cp ./artifacts/zkcensus.pkey ./wasm/plonk/zkcensus.pkey
	@echo "copying compatible wasm_exec.js"
	@cp ./artifacts/wasm_exec_go.js ./examples/web/wasm_exec.js
	@echo "compilling the prover for go-wasm"
	@make compile-prover-go-plonk
	@echo "removing copied artifacts"
	@rm ./wasm/plonk/zkcensus.ccs
	@rm ./wasm/plonk/zkcensus.srs
	@rm ./wasm/plonk/zkcensus.pkey
	@rm -f examples/web/index.html
	@ln -s index_plonk.html examples/web/index.html
	@go run examples/web/main.go

run-tinygo-web-example-plonk:
	@echo "compilling circuit and genering artifacts"
	@make compile-circuit-plonk
	@echo "copying artifacts"
	@cp ./artifacts/zkcensus.ccs ./wasm/plonk/zkcensus.ccs
	@cp ./artifacts/zkcensus.srs ./wasm/plonk/zkcensus.srs
	@cp ./artifacts/zkcensus.pkey ./wasm/plonk/zkcensus.pkey
	@echo "compilling the prover for tinygo"
	@make compile-prover-tinygo-plonk
	@echo "removing copied artifacts"
	@rm ./wasm/plonk/zkcensus.ccs
	@rm ./wasm/plonk/zkcensus.srs
	@rm ./wasm/plonk/zkcensus.pkey
	@echo "copying compatible wasm_exec.js"
	@cp ./artifacts/wasm_exec_tinygo.js ./examples/web/wasm_exec.js
	@rm -f examples/web/index.html
	@ln -s index_plonk.html examples/web/index.html
	@go run examples/web/main.go

run-tinygo-web-example-g16:
	@echo "compilling circuit and genering artifacts for groth16"
	@make compile-circuit-g16
	@echo "copying artifacts"
	@cp ./artifacts/g16_zkcensus.ccs ./wasm/g16/zkcensus.ccs
	@cp ./artifacts/g16_zkcensus.pkey ./wasm/g16/zkcensus.pkey
	@echo "compilling the prover for tinygo"
	@make compile-prover-tinygo-g16
	@echo "removing copied artifacts"
	@rm ./wasm/g16/zkcensus.ccs
	@rm ./wasm/g16/zkcensus.pkey
	@rm -f examples/web/index.html
	@echo "copying compatible wasm_exec.js"
	@cp ./artifacts/wasm_exec_tinygo.js ./examples/web/wasm_exec.js
	@ln -s index_g16.html examples/web/index.html
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

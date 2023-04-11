//go:build tinygo
// +build tinygo

package main

import (
	_ "embed"
	"fmt"
	"gnark-prover-tinygo/prover"
	"syscall/js"
)

//go:embed zkcensus.ccs
var circuit []byte

//go:embed zkcensus.pkey
var epkey []byte

func main() {
	c := make(chan int)
	js.Global().Set("generateProof", js.FuncOf(jsGenerateProof))
	<-c
}

func jsGenerateProof(this js.Value, args []js.Value) interface{} {
	bwitness := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(bwitness, args[0])
	fmt.Println("Calling function GenerateProof")
	if _, _, err := prover.GenerateProofGroth16(circuit, epkey, bwitness); err != nil {
		return err.Error()
	}
	return nil
}

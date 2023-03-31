//go:build js && wasm
// +build js,wasm

package main

import (
	"gnark-prover-tinygo/internal/circuit/plonk"
	"log"
	"syscall/js"
)

func main() {
	c := make(chan int)
	js.Global().Set("generateProof", js.FuncOf(jsGenerateProof))
	<-c
}

func jsGenerateProof(this js.Value, args []js.Value) interface{} {
	// var bccs, bsrs, witness []byte
	bccs := make([]byte, args[0].Get("length").Int())
	bsrs := make([]byte, args[0].Get("length").Int())
	bwitness := make([]byte, args[0].Get("length").Int())

	js.CopyBytesToGo(bccs, args[0])
	js.CopyBytesToGo(bsrs, args[1])
	js.CopyBytesToGo(bwitness, args[2])

	if _, _, _, err := plonk.GenerateProof(bccs, bsrs, bwitness); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

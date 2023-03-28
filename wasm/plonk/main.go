//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
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

	vk, proof, pubWitness, err := plonk.GenerateProof(bccs, bsrs, bwitness)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(vk, proof, pubWitness)
	}
	return true
}

func jsVerifyProof(bsrs, bvk, bproof, bpubwitness []byte) error {
	return plonk.VerifyProof(bsrs, bvk, bproof, bpubwitness)
}

//export GenerateProof
func GenerateProof(bccs, bsrs, inputs []byte) ([]byte, []byte, []byte, error) {
	return plonk.GenerateProof(bccs, bsrs, inputs)
}

//export VerifyProof
func VerifyProof(bsrs, bvk, bproof, bpubwitness []byte) error {
	return plonk.VerifyProof(bsrs, bvk, bproof, bpubwitness)
}
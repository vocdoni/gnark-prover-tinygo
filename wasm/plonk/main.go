//go:build tinygo
// +build tinygo

package main

import (
	_ "embed"
	"fmt"
	"gnark-prover-tinygo/prover"
	"syscall/js"
	"time"
)

//go:embed zkcensus.ccs
var eccs []byte

//go:embed zkcensus.srs
var esrs []byte

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
	startTime := time.Now()
	if _, _, err := prover.GenerateProofPlonk(eccs, esrs, epkey, bwitness); err != nil {
		fmt.Println("Error calling function GenerateProof", err.Error())
		return 0
	}
	elapsedTime := int(time.Now().Sub(startTime).Seconds() * 1000)
	return elapsedTime
}

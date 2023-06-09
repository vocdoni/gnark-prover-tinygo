//go:build tinygo
// +build tinygo

package main

import (
	_ "embed"
	"fmt"
	"gnark-prover-tinygo/prover"
)

//go:embed zkcensus.ccs
var circuit []byte

//go:embed zkcensus.pkey
var epkey []byte

func main() {
	c := make(chan interface{})
	<-c
}

//export generateProof
func GenerateProof(bwitness []byte) interface{} {
	if _, _, err := prover.GenerateProofGroth16(circuit, epkey, bwitness); err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return nil
}

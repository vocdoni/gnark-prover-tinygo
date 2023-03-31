//go:build tinygo
// +build tinygo

package main

import "gnark-prover-tinygo/internal/circuit/groth16"

func main() {}

//export GenerateProof
func GenerateProof(bccs, bsrs, inputs []byte) ([]byte, []byte, []byte, error) {
	return groth16.GenerateProof(bccs, bsrs, inputs)
}

//go:build tinygo
// +build tinygo

package main

import "gnark-prover-tinygo/internal/circuit/plonk"

func main() {}

//export GenerateProof
func GenerateProof(bccs, bsrs, inputs []byte) ([]byte, []byte, []byte, error) {
	return plonk.GenerateProof(bccs, bsrs, inputs)
}

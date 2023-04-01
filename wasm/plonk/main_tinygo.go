//go:build tinygo
// +build tinygo

package main

import "gnark-prover-tinygo/internal/circuit/plonk"

func main() {}

//export GenerateProof
func GenerateProof(bccs, bsrs, bpkey, inputs []byte) ([]byte, []byte, error) {
	return plonk.GenerateProof(bccs, bsrs, bpkey, inputs)
}

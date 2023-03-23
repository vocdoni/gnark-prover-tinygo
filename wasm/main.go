package main

import "gnark-test/internal/circuit"

func main() {}

//export GenerateProof
func GenerateProof(bccs, bsrs, inputs []byte) ([]byte, []byte, []byte, error) {
	return circuit.GenerateProof(bccs, bsrs, inputs)
}

//export VerifyProof
func VerifyProof(bsrs, bvk, bproof, bpubwitness []byte) error {
	return circuit.VerifyProof(bsrs, bvk, bproof, bpubwitness)
}

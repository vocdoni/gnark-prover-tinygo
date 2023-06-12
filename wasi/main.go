//go:build tinygo
// +build tinygo

package main

import (
	_ "embed"
	"fmt"
	"gnark-prover-tinygo/prover"
	"unsafe"
)

//go:embed zkcensus.ccs
var circuit []byte

//go:embed zkcensus.pkey
var epkey []byte

//go:embed witness.bin
var witness []byte

func main() {
	getProof()
}

//export start
func start() {
	fmt.Println("start!!!")
	getProof()
}

//export getProof
func getProof() {
	fmt.Println("executing generateProof")
	fmt.Printf("witness size %d\n", len(witness))
	if _, _, err := prover.GenerateProofGroth16(circuit, epkey, witness); err != nil {
		fmt.Println(err)
	}
}

//export generateProof
func GenerateProof(witnessPtr *byte, witnessLen int) {
	fmt.Println("witness len", witnessLen)
	fmt.Println("proving key len", len(epkey))
	witnessSlice := (*[1 << 28]byte)(unsafe.Pointer(witnessPtr))[:witnessLen:witnessLen]
	if _, _, err := prover.GenerateProofGroth16(circuit, epkey, witnessSlice); err != nil {
		fmt.Println(err)
	}
}

//export alloc
func alloc(size int) unsafe.Pointer {
	return unsafe.Pointer(&make([]byte, size)[0])
}

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

func main() {
}

//export start
func start() {
	fmt.Println("start!!!")
}

//export getProof
func getProof(bwitness []byte) {
	fmt.Println("executing generateProof")
	fmt.Printf("witness size %d\n", len(bwitness))
	if _, _, err := prover.GenerateProofGroth16(circuit, epkey, bwitness); err != nil {
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

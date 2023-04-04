package main

import (
	"flag"
	"fmt"
	"os"

	cp "gnark-prover-tinygo/prover"
)

func main() {
	fdcircuit := flag.String("circuit", "", "circuit file")
	fdsrs := flag.String("srs", "", "srs file")
	fdpkey := flag.String("pkey", "", "proving key file")
	fdwitness := flag.String("witness", "", "witness file")

	flag.Parse()

	// Read the files into byte slices and call the generateProof function
	fmt.Println("reading circuit file: ", *fdcircuit)
	bccs, err := os.ReadFile(*fdcircuit)
	if err != nil {
		panic(err)
	}

	fmt.Println("reading srs file: ", *fdsrs)
	bsrs, err := os.ReadFile(*fdsrs)
	if err != nil {
		panic(err)
	}

	fmt.Println("reading proving key file: ", *fdpkey)
	bpkey, err := os.ReadFile(*fdpkey)
	if err != nil {
		panic(err)
	}

	fmt.Println("reading witness file: ", *fdwitness)
	bwitness, err := os.ReadFile(*fdwitness)
	if err != nil {
		panic(err)
	}

	fmt.Println("calling generateProof function")
	proof, publicWitness, err := cp.GenerateProof(bccs, bsrs, bpkey, bwitness)
	if err != nil {
		panic(err)
	}
	fmt.Printf("proof: %x\n", proof)
	fmt.Printf("public witness: %x\n", publicWitness)
}

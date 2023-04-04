package main

import (
	"flag"
	"fmt"
	"gnark-prover-tinygo/circuits/zkcensus"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
)

var ccsOutput = flag.String("ccs", "./artifacts/zkcensus.ccs", "Output file to the encoded output Gnark Circuit Constrain System, result of circuit compilation")
var srsOutput = flag.String("srs", "./artifacts/zkcensus.srs", "Output file to the encoded output Gnark KZG polynomial commitment, result of circuit compilation")
var pKeyOutput = flag.String("pkey", "./artifacts/zkcensus.pkey", "Circuit proving key")
var vKeyOutput = flag.String("vkey", "./artifacts/zkcensus.vkey", "Circuit verifying key")
var witnessOutput = flag.String("witness", "./artifacts/zkcensus.witness", "Circuit witness")

func main() {
	flag.Parse()

	ccs, srs, pKey, vKey, err := compilePlonk()
	if err != nil {
		log.Fatal(err)
	}
	if err := savePlonk(ccs, srs, pKey, vKey, *ccsOutput, *srsOutput, *pKeyOutput, *vKeyOutput); err != nil {
		log.Fatal(err)
	}

	input, _ := zkcensus.ZkCensusInputs(160, 100)
	witness, _ := frontend.NewWitness(&input, ecc.BN254.ScalarField())
	fdWitness, err := os.Create(*witnessOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer fdWitness.Close()
	if _, err := witness.WriteTo(fdWitness); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", input.String())
}

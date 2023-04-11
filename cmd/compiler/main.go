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

var (
	ccsOutput     = flag.String("ccs", "./artifacts/zkcensus.ccs", "Output file to the encoded output Gnark Circuit Constrain System, result of circuit compilation")
	srsOutput     = flag.String("srs", "./artifacts/zkcensus.srs", "Output file to the encoded output Gnark KZG polynomial commitment, result of circuit compilation")
	pKeyOutput    = flag.String("pkey", "./artifacts/zkcensus.pkey", "Circuit proving key")
	vKeyOutput    = flag.String("vkey", "./artifacts/zkcensus.vkey", "Circuit verifying key")
	witnessOutput = flag.String("witness", "./artifacts/zkcensus.witness", "Circuit witness")

	g16circuitOutput = flag.String("g16circuit", "./artifacts/g16_zkcensus.ccs", "Output file to the encoded output Gnark Circuit Constrain System, result of circuit compilation")
	g16pKeyOutput    = flag.String("g16pkey", "./artifacts/g16_zkcensus.pkey", "Circuit proving key")
	g16vKeyOutput    = flag.String("g16vkey", "./artifacts/g16_zkcensus.vkey", "Circuit verifying key")

	protocol = flag.String("protocol", "plonk", "Protocol to use: plonk or groth16")
)

func main() {
	flag.Parse()
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

	switch *protocol {
	case "plonk":
		ccs, srs, pKey, vKey, err := compilePlonk()
		if err != nil {
			log.Fatal(err)
		}
		if err := savePlonk(ccs, srs, pKey, vKey, *ccsOutput, *srsOutput, *pKeyOutput, *vKeyOutput); err != nil {
			log.Fatal(err)
		}
	case "groth16":
		g16circuit, g16pKey, g16vKey, err := compileGroth16()
		if err != nil {
			log.Fatal(err)
		}
		if err := saveGroth16(g16circuit, g16pKey, g16vKey, *g16circuitOutput, *g16pKeyOutput, *g16vKeyOutput); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%s", input.String())
}

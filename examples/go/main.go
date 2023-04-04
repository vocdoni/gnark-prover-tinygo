package main

import (
	"log"
	"os"
	"time"

	"gnark-prover-tinygo/prover"
)

func main() {
	ccs, err := os.ReadFile("./artifacts/zkcensus.ccs")
	if err != nil {
		panic(err)
	}
	srs, _ := os.ReadFile("./artifacts/zkcensus.srs")
	pkey, err := os.ReadFile("./artifacts/zkcensus.pkey")
	if err != nil {
		panic(err)
	}
	witness, err := os.ReadFile("./artifacts/zkcensus.witness")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	proof, pubWitness, err := prover.GenerateProof(ccs, srs, pkey, witness)
	if err != nil {
		panic(err)
	}
	log.Println(proof, pubWitness)
	log.Println("Took", time.Since(start))
}

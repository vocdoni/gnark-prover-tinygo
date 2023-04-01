package main

import (
	"gnark-prover-tinygo/internal/circuit/groth16"
	"gnark-prover-tinygo/internal/circuit/plonk"
	"log"
	"os"
	"time"
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
	proof, pubWitness, err := plonk.GenerateProof(ccs, srs, pkey, witness)
	if err == nil {
		log.Println(proof, pubWitness)
		log.Println("Took", time.Since(start))
		return
	}

	start = time.Now()
	proof, pubWitness, err = groth16.GenerateProof(ccs, pkey, witness)
	if err != nil {
		panic(err)
	}
	log.Println(proof, pubWitness)
	log.Println("Took", time.Since(start))
}

package main

import (
	"gnark-prover-tinygo/internal/circuit/groth16"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ccs, err := os.ReadFile("./artifacts/zkcensus.ccs")
	if err != nil {
		panic(err)
	}
	srs, err := os.ReadFile("./artifacts/zkcensus.srs")
	if err != nil {
		panic(err)
	}
	witness, err := os.ReadFile("./artifacts/witness")
	if err != nil {
		panic(err)
	}

	vk, proof, pubWitness, err := groth16.GenerateProof(ccs, srs, witness)
	if err != nil {
		panic(err)
	}

	log.Println(vk, proof, pubWitness)
	log.Println("Took", time.Since(start))
}

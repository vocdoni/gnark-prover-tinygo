package main

import (
	"flag"
	"gnark-prover-tinygo/internal/circuit/groth16"
	"gnark-prover-tinygo/internal/circuit/plonk"
	"log"
	"os"
	"time"
)

var zkBackend = flag.String("backend", "groth16", "backend circuit ('groth16' or 'plonk')")

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
	var vk, proof, pubWitness []byte
	switch *zkBackend {
	case "plonk":
		vk, proof, pubWitness, err = plonk.GenerateProof(ccs, srs, witness)
		if err != nil {
			panic(err)
		}
	case "groth16":
		vk, proof, pubWitness, err = groth16.GenerateProof(ccs, srs, witness)
		if err != nil {
			panic(err)
		}
	}
	log.Println(vk, proof, pubWitness)
	log.Println("Took", time.Since(start))
}

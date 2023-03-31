package main

import (
	"log"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/test"
)

func TestMiMC(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit TestMiMCCircuit
	assignment := TestMiMCCircuit{
		Input: 1995,
	}

	start := time.Now()
	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
	log.Println("MiMC Plonk took (s):", time.Since(start))

	start = time.Now()
	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))
	log.Println("MiMC Groth16 took (s):", time.Since(start))
}

func TestPoseidon(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit TestPoseidonCircuit
	assignment := TestPoseidonCircuit{
		Input: 1995,
	}

	start := time.Now()
	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
	log.Println("Poseidon Plonk took (s):", time.Since(start))

	start = time.Now()
	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))
	log.Println("Poseidon Groth16 took (s):", time.Since(start))
}

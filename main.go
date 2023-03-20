package main

import (
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/permutation/keccakf"
)

// CubicCircuit defines a simple circuit
// x**3 + x + 5 == y
type CubicCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *CubicCircuit) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
	kdata := [25]frontend.Variable{}
	for i := 0; i < 25; i++ {
		kdata[i] = x3
	}
	keccakf.Permute(api, kdata)
	return nil
}

func writeCircuit(ccs constraint.ConstraintSystem) {
	// compiles our circuit into a R1CS
	// write the circuit to a file
	f, err := os.Create("circuit.r1cs")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := ccs.WriteTo(f); err != nil {
		panic(err)
	}
}

func main() {
	fcircuit, err := os.Open("circuit.r1cs")
	if err != nil {
		panic(err)
	}
	defer fcircuit.Close()
	ccs := plonk.NewCS(ecc.BN254)
	if _, err := ccs.ReadFrom(fcircuit); err != nil {
		panic(err)
	}

	// load the KZG srs from file
	srs := kzg.NewSRS(ecc.BN254)
	fkzg, err := os.Open("kzg.bin")
	if err != nil {
		panic(err)
	}
	defer fkzg.Close()
	if _, err := srs.ReadFrom(fkzg); err != nil {
		panic(err)
	}

	// Generate the proving and verification keys.
	pk, _, err := plonk.Setup(ccs, srs)
	if err != nil {
		log.Fatal(err)
	}

	// witness definition
	assignment := CubicCircuit{X: 3, Y: 35}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())

	// Prove & Verify
	proof, err := plonk.Prove(ccs, pk, witness)
	if err != nil {
		panic(err)
	}

	proof.WriteTo(os.Stdout)
}

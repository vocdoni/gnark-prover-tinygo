package main

import (
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
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
	cs := groth16.NewCS(ecc.BN254)
	if _, err := cs.ReadFrom(fcircuit); err != nil {
		panic(err)
	}

	fpk, err := os.Open("proving_key.bin")
	if err != nil {
		panic(err)
	}
	defer fpk.Close()
	pk := groth16.NewProvingKey(ecc.BN254)
	if _, err := pk.ReadFrom(fpk); err != nil {
		panic(err)
	}

	// witness definition
	assignment := CubicCircuit{X: 3, Y: 35}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())

	// groth16: Prove & Verify
	proof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		panic(err)
	}
	proof.WriteTo(os.Stdout)
}

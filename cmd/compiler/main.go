package main

import (
	"gnark-test/circuits/zkcensus"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
)

func compileZkCensus() (constraint.ConstraintSystem, error) {
	var c zkcensus.ZkCensusCircuit
	return frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &c)
}

func saveCompiledCircuit(ccs constraint.ConstraintSystem, dst string) error {
	fd, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = ccs.WriteTo(fd)
	return err
}

func main() {
	cs, err := compileZkCensus()
	if err != nil {
		log.Fatalln(err)
	}

	if err := saveCompiledCircuit(cs, "./artifacts/zkcensus.r1cs"); err != nil {
		return
	}
}

package main

import (
	"gnark-prover-tinygo/circuits/zkcensus"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
)

func compileGroth16() (constraint.ConstraintSystem, kzg.SRS, error) {
	var c zkcensus.ZkCensusCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &c)
	if err != nil {
		return nil, nil, err
	}

	srs, err := test.NewKZGSRS(ccs)
	if err != nil {
		return nil, nil, err
	}
	return ccs, srs, nil
}

func saveGroth16(ccs constraint.ConstraintSystem, srs kzg.SRS, ccsDst, srsDst string) error {
	fdSRS, err := os.Create(srsDst)
	if err != nil {
		return err
	}
	defer fdSRS.Close()
	if _, err = srs.WriteTo(fdSRS); err != nil {
		return err
	}

	fdCCS, err := os.Create(ccsDst)
	if err != nil {
		return err
	}
	defer fdCCS.Close()
	_, err = ccs.WriteTo(fdCCS)
	return err
}

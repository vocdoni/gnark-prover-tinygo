package main

import (
	"gnark-prover-tinygo/circuits/zkcensus"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
)

func compilePlonk() (constraint.ConstraintSystem, kzg.SRS, error) {
	var c zkcensus.ZkCensusCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &c)
	if err != nil {
		return nil, nil, err
	}

	srs, err := test.NewKZGSRS(ccs)
	if err != nil {
		return nil, nil, err
	}
	return ccs, srs, nil
}

func savePlonk(ccs constraint.ConstraintSystem, srs kzg.SRS, ccsDst, srsDst string) error {
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

	_scs := ccs.(*cs.SparseR1CS)
	_, err = _scs.WriteTo(fdCCS)
	return err
}

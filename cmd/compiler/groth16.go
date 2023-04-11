package main

import (
	"gnark-prover-tinygo/circuits/zkcensus"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func compileGroth16() (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey, error) {
	var c zkcensus.ZkCensusCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &c)
	if err != nil {
		return nil, nil, nil, err
	}

	provingKey, verifyingKey, err := groth16.Setup(ccs)
	if err != nil {
		return nil, nil, nil, err
	}

	return ccs, provingKey, verifyingKey, nil
}

func saveGroth16(ccs constraint.ConstraintSystem, provingKey groth16.ProvingKey,
	verifyingKey groth16.VerifyingKey, ccsDst, pKeyDst, vKeyDst string) error {
	fdCCS, err := os.Create(ccsDst)
	if err != nil {
		return err
	}
	defer fdCCS.Close()
	_r1cs := ccs.(*cs.R1CS)
	if _, err := _r1cs.WriteTo(fdCCS); err != nil {
		return err
	}

	fdPKey, err := os.Create(pKeyDst)
	if err != nil {
		return err
	}
	defer fdPKey.Close()
	if _, err := provingKey.WriteTo(fdPKey); err != nil {
		return err
	}

	fdVKey, err := os.Create(vKeyDst)
	if err != nil {
		return err
	}
	defer fdVKey.Close()
	if _, err := verifyingKey.WriteTo(fdVKey); err != nil {
		return err
	}
	return nil
}

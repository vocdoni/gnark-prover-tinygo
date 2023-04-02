package main

import (
	"crypto/rand"
	"fmt"
	"gnark-prover-tinygo/circuits/zkcensus"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/kzg"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/vocdoni/gnark-wasm-prover/encoder"
	"github.com/vocdoni/gnark-wasm-prover/utils"
)

const srsKZGsize = (1 << 14) + 3

func compilePlonk() (constraint.ConstraintSystem, *kzg.SRS, plonk.ProvingKey, plonk.VerifyingKey, error) {
	var c zkcensus.ZkCensusCircuit

	// Compile circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &c)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Generate KZG file (trusted setup)
	curveID := utils.FieldToCurve(ccs.Field())
	alpha, err := rand.Int(rand.Reader, curveID.ScalarField())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// calculate the size for the KZG
	nbConstraints := ccs.GetNbConstraints()
	sizeSystem := nbConstraints + ccs.GetNbPublicVariables()
	kzgSize := ecc.NextPowerOfTwo(uint64(sizeSystem)) + 3

	fmt.Println("Generating KZG SRS for curve", curveID.String(), "with size", kzgSize)
	kzg, err := kzg.NewSRS(kzgSize, alpha)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Generate proving and verifying keys
	fmt.Println("Generating proving and verifying keys")
	provingKey, verifyingKey, err := plonk.Setup(ccs, kzg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return ccs, kzg, provingKey, verifyingKey, nil
}

func savePlonk(ccs constraint.ConstraintSystem, kzg *kzg.SRS, provingKey plonk.ProvingKey, verifyingKey plonk.VerifyingKey, ccsDst, srsDst, pKeyDst, vKeyDst string) error {
	fdSRS, err := os.Create(srsDst)
	if err != nil {
		return err
	}
	defer fdSRS.Close()
	// KZG has its own binary encoding
	n, err := kzg.WriteTo(fdSRS)
	if err != nil {
		return err
	}
	println("KZG size:", n)

	fdCCS, err := os.Create(ccsDst)
	if err != nil {
		return err
	}
	defer fdCCS.Close()
	_scs := ccs.(*cs.SparseR1CS)
	n, err = encoder.EncodeToGob(fdCCS, _scs)
	if err != nil {
		return err
	}
	println("CCS Gob size:", n)

	fdPKey, err := os.Create(pKeyDst)
	if err != nil {
		return err
	}
	defer fdPKey.Close()
	n, err = provingKey.WriteTo(fdPKey)
	if err != nil {
		return err
	}
	println("ProvingKey size:", n)

	fdVKey, err := os.Create(vKeyDst)
	if err != nil {
		return err
	}
	defer fdVKey.Close()
	n, err = verifyingKey.WriteTo(fdVKey)
	if err != nil {
		return err
	}
	println("VerifyingKey size:", n)

	return nil
}

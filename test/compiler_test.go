package main

import (
	"log"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
)

func compileGroth16(c frontend.Circuit) (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey, error) {
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, c)
	if err != nil {
		return nil, nil, nil, err
	}

	provingKey, verifyingKey, err := groth16.Setup(ccs)
	if err != nil {
		return nil, nil, nil, err
	}

	return ccs, provingKey, verifyingKey, nil
}

func saveGroth16(ccs constraint.ConstraintSystem, provingKey groth16.ProvingKey, verifyingKey groth16.VerifyingKey, ccsDst, pKeyDst, vKeyDst string) error {
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

func compilePlonk(c frontend.Circuit) (constraint.ConstraintSystem, kzg.SRS, plonk.ProvingKey, plonk.VerifyingKey, error) {
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, c)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	srs, err := test.NewKZGSRS(ccs)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	provingKey, verifyingKey, err := plonk.Setup(ccs, srs)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return ccs, srs, provingKey, verifyingKey, nil
}

func savePlonk(ccs constraint.ConstraintSystem, srs kzg.SRS, provingKey plonk.ProvingKey, verifyingKey plonk.VerifyingKey, ccsDst, srsDst, pKeyDst, vKeyDst string) error {
	fdCCS, err := os.Create(ccsDst)
	if err != nil {
		return err
	}
	defer fdCCS.Close()
	_scs := ccs.(*cs.SparseR1CS)
	if _, err := _scs.WriteTo(fdCCS); err != nil {
		return err
	}

	fdSRS, err := os.Create(srsDst)
	if err != nil {
		return err
	}
	defer fdSRS.Close()
	if _, err = srs.WriteTo(fdSRS); err != nil {
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
	if _, err := provingKey.WriteTo(fdVKey); err != nil {
		return err
	}

	return nil
}

func TestCompileGroth16(t *testing.T) {
	if ccs, pkey, vkey, err := compileGroth16(&TestMiMCCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := saveGroth16(ccs, pkey, vkey, "./artifacts/mimc_groth16.ccs", "./artifacts/mimc_groth16.pkey", "./artifacts/mimc_groth16.vkey"); err != nil {
		log.Fatal(err)
	}

	if ccs, pkey, vkey, err := compileGroth16(&TestPoseidonCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := saveGroth16(ccs, pkey, vkey, "./artifacts/poseidon_groth16.ccs", "./artifacts/poseidon_groth16.pkey", "./artifacts/poseidon_groth16.vkey"); err != nil {
		log.Fatal(err)
	}
}

func TestCompilePlonk(t *testing.T) {
	if ccs, srs, pkey, vkey, err := compilePlonk(&TestMiMCCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := savePlonk(ccs, srs, pkey, vkey, "./artifacts/mimc_plonk.ccs", "./artifacts/mimc_plonk.srs", "./artifacts/mimc_plonk.pkey", "./artifacts/mimc_plonk.vkey"); err != nil {
		log.Fatal(err)
	}

	if ccs, srs, pkey, vkey, err := compilePlonk(&TestPoseidonCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := savePlonk(ccs, srs, pkey, vkey, "./artifacts/poseidon_plonk.ccs", "./artifacts/poseidon_plonk.srs", "./artifacts/poseidon_plonk.pkey", "./artifacts/poseidon_plonk.vkey"); err != nil {
		log.Fatal(err)
	}
}

func TestGenerateWitness(t *testing.T) {
	witnessMimc, _ := frontend.NewWitness(&TestMiMCCircuit{
		Input: 1000,
	}, ecc.BN254.ScalarField())
	f, err := os.Create("./artifacts/witness_mimc")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = witnessMimc.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}

	witnessPoseidon, _ := frontend.NewWitness(&TestPoseidonCircuit{
		Input: 1000,
	}, ecc.BN254.ScalarField())

	f, err = os.Create("./artifacts/witness_poseidon")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = witnessPoseidon.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
)

func compileGroth16(c frontend.Circuit) (constraint.ConstraintSystem, kzg.SRS, error) {
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, c)
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
	_r1cs := ccs.(*cs.R1CS)
	_, err = _r1cs.WriteTo(fdCCS)
	return err
}

func compilePlonk(c frontend.Circuit) (constraint.ConstraintSystem, kzg.SRS, error) {
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, c)
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

func TestCompileGroth16(t *testing.T) {
	if ccs, srs, err := compileGroth16(&TestMiMCCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := saveGroth16(ccs, srs, "./artifacts/mimc_groth16.ccs", "./artifacts/mimc_groth16.srs"); err != nil {
		log.Fatal(err)
	}

	if ccs, srs, err := compileGroth16(&TestPoseidonCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := saveGroth16(ccs, srs, "./artifacts/poseidon_groth16.ccs", "./artifacts/poseidon_groth16.srs"); err != nil {
		log.Fatal(err)
	}
}

func TestCompilePlonk(t *testing.T) {
	if ccs, srs, err := compilePlonk(&TestMiMCCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := savePlonk(ccs, srs, "./artifacts/mimc_plonk.ccs", "./artifacts/mimc_plonk.srs"); err != nil {
		log.Fatal(err)
	}

	if ccs, srs, err := compilePlonk(&TestPoseidonCircuit{}); err != nil {
		log.Fatal(err)
	} else if err := savePlonk(ccs, srs, "./artifacts/poseidon_plonk.ccs", "./artifacts/poseidon_plonk.srs"); err != nil {
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

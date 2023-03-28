package groth16

import (
	"bytes"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"

	// This import fixes the issue that raises when a prover tries to generate a proof
	// of a serialized circuit. Check more information here:
	//   - https://github.com/ConsenSys/gnark/issues/600
	//   - https://github.com/phated/gnark-browser/blob/2446c65e89156f1a04163724a89e5dcb7e4c4886/README.md#solution-hint-registration
	_ "github.com/consensys/gnark/std/math/bits"
)

// GenerateProof sets up the circuit with the constrain system and the srs files
// provided and generates the proof for the JSON encoded inputs (witness). It
// returns the verification key, the proof and the public witness, all of this
// outputs will be encoded as JSON. If something fails, it returns an error.
func GenerateProof(bccs, bsrs, inputs []byte) ([]byte, []byte, []byte, error) {
	// Read and initialize circuit CS
	ccs := groth16.NewCS(ecc.BN254)
	if _, err := ccs.ReadFrom(bytes.NewReader(bccs)); err != nil {
		return nil, nil, nil, fmt.Errorf("error reading circuit cs: %w", err)
	}
	// Read and initialize the witness
	cWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error initializing witness: %w", err)
	}
	if _, err := cWitness.ReadFrom(bytes.NewReader(inputs)); err != nil {
		return nil, nil, nil, fmt.Errorf("error reading witness: %w", err)
	}
	// Get proving and verifiying keys
	provingKey, verifyingKey, err := groth16.Setup(ccs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error generating plonk keys: %w", err)
	}

	// Generate the proof
	proof, err := groth16.Prove(ccs, provingKey, cWitness)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error generating proof: %w", err)
	}
	proofBuff := bytes.Buffer{}
	if _, err := proof.WriteTo(&proofBuff); err != nil {
		return nil, nil, nil, fmt.Errorf("error encoding proof: %w", err)
	}
	// Get public witness part and encode it
	publicWitness, err := cWitness.Public()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error generating public witness: %w", err)
	}
	publicWitnessBuff := bytes.Buffer{}
	if _, err := publicWitness.WriteTo(&publicWitnessBuff); err != nil {
		return nil, nil, nil, fmt.Errorf("error encoding public witness: %w", err)
	}
	// Encode verifiying key
	verifyingKeyBuff := bytes.Buffer{}
	if _, err := verifyingKey.WriteTo(&verifyingKeyBuff); err != nil {
		return nil, nil, nil, fmt.Errorf("error encoding verifiying key: %w", err)
	}
	return verifyingKeyBuff.Bytes(),
		proofBuff.Bytes(),
		publicWitnessBuff.Bytes(),
		nil
}

// VerifyProof verifies the proof provided using the verifiying key and the
// public witness also provided. It returns an error if something fails during
// inputs parsing or proof verification.
func VerifyProof(bsrs, bvk, bproof, bpubwitness []byte) error {
	// Parse the verifiying key
	verifiyingKey := groth16.NewVerifyingKey(ecc.BN254)
	if _, err := verifiyingKey.ReadFrom(bytes.NewBuffer(bvk)); err != nil {
		return fmt.Errorf("error reading verifiying key: %w", err)
	}
	// Parse the proof
	proof := groth16.NewProof(ecc.BN254)
	if _, err := proof.ReadFrom(bytes.NewBuffer(bproof)); err != nil {
		return fmt.Errorf("error reading proof: %w", err)
	}
	// Parse the public witness
	pubWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return fmt.Errorf("error initializing public witness: %w", err)
	}
	if _, err := pubWitness.ReadFrom(bytes.NewReader(bpubwitness)); err != nil {
		return fmt.Errorf("error reading public witness: %w", err)
	}
	// Return the result of the verification process
	return groth16.Verify(proof, verifiyingKey, pubWitness)
}

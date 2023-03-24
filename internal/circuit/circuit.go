package circuit

import (
	"bytes"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/pkg/errors"
)

// GenerateProof sets up the circuit with the constrain system and the srs files
// provided and generates the proof for the JSON encoded inputs (witness). It
// returns the verification key, the proof and the public witness, all of this
// outputs will be encoded as JSON. If something fails, it returns an error.
func GenerateProof(bccs, bsrs, inputs []byte) ([]byte, []byte, []byte, error) {
	// Read and initialize circuit CS
	ccs := plonk.NewCS(ecc.BN254)
	if _, err := ccs.ReadFrom(bytes.NewReader(bccs)); err != nil {
		return nil, nil, nil, errors.Wrap(err, "error reading circuit cs")
	}
	// Read and initialize SSR
	srs := kzg.NewSRS(ecc.BN254)
	if _, err := srs.ReadFrom(bytes.NewReader(bsrs)); err != nil {
		return nil, nil, nil, errors.Wrap(err, "error reading plonk srs")
	}
	// Read and initialize the witness
	cWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error initializing witness")
	}
	if _, err := cWitness.ReadFrom(bytes.NewReader(inputs)); err != nil {
		return nil, nil, nil, errors.Wrap(err, "error reading witness")
	}
	// Get proving and verifiying keys
	provingKey, verifyingKey, err := plonk.Setup(ccs, srs)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error generating plonk keys")
	}

	// var c zkcensus.ZkCensusCircuit
	// r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &c)

	// Generate the proof
	proof, err := plonk.Prove(ccs, provingKey, cWitness)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error generating proof")
	}
	proofBuff := bytes.Buffer{}
	if _, err := proof.WriteTo(&proofBuff); err != nil {
		return nil, nil, nil, errors.Wrap(err, "error encoding proof")
	}
	// Get public witness part and encode it
	publicWitness, err := cWitness.Public()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error generating public witness")
	}
	publicWitnessBuff := bytes.Buffer{}
	if _, err := publicWitness.WriteTo(&publicWitnessBuff); err != nil {
		return nil, nil, nil, errors.Wrap(err, "error encoding public witness")
	}
	// Encode verifiying key
	verifyingKeyBuff := bytes.Buffer{}
	if _, err := verifyingKey.WriteTo(&verifyingKeyBuff); err != nil {
		return nil, nil, nil, errors.Wrap(err, "error encoding verifiying key")
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
	verifiyingKey := plonk.NewVerifyingKey(ecc.BN254)
	if _, err := verifiyingKey.ReadFrom(bytes.NewBuffer(bvk)); err != nil {
		return errors.Wrap(err, "error reading verifiying key")
	}
	// Read and initialize SSR
	srs := kzg.NewSRS(ecc.BN254)
	if _, err := srs.ReadFrom(bytes.NewReader(bsrs)); err != nil {
		return errors.Wrap(err, "error reading plonk srs")
	}
	if err := verifiyingKey.InitKZG(srs); err != nil {
		return errors.Wrap(err, "error initializing srs verifiying key")
	}
	// Parse the proof
	proof := plonk.NewProof(ecc.BN254)
	if _, err := proof.ReadFrom(bytes.NewBuffer(bproof)); err != nil {
		return errors.Wrap(err, "error reading proof")
	}
	// Parse the public witness
	pubWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return errors.Wrap(err, "error initializing public witness")
	}
	if _, err := pubWitness.ReadFrom(bytes.NewReader(bpubwitness)); err != nil {
		return errors.Wrap(err, "error reading public witness")
	}
	// Return the result of the verification process
	return plonk.Verify(proof, verifiyingKey, pubWitness)
}

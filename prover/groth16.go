package prover

import (
	"bytes"
	"fmt"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/std"

	// This import fixes the issue that raises when a prover tries to generate a proof
	// of a serialized circuit. Check more information here:
	//   - https://github.com/ConsenSys/gnark/issues/600
	//   - https://github.com/phated/gnark-browser/blob/2446c65e89156f1a04163724a89e5dcb7e4c4886/README.md#solution-hint-registration
	_ "github.com/consensys/gnark/std/math/bits"
)

// GenerateProofGroth16 sets up the circuit with the constrain system and the srs files
// provided and generates the proof for the JSON encoded inputs (witness). It
// returns the verification key, the proof and the public witness, all of this
// outputs will be encoded as JSON. If something fails, it returns an error.
func GenerateProofGroth16(bccs, bpkey, inputs []byte) ([]byte, []byte, error) {
	step := time.Now()
	// Read and initialize circuit CS
	ccs := groth16.NewCS(ecc.BN254)
	if _, err := ccs.ReadFrom(bytes.NewReader(bccs)); err != nil {
		fmt.Println("error reading circuit cs: ", err)
		return nil, nil, fmt.Errorf("error reading circuit cs: %w", err)
	}
	fmt.Println("ccs loaded, took (s):", time.Since(step))
	step = time.Now()
	// Read proving key
	provingKey := groth16.NewProvingKey(ecc.BN254)
	if _, err := provingKey.UnsafeReadFrom(bytes.NewReader(bpkey)); err != nil {
		fmt.Println("error reading circuit pkey: ", err)
		return nil, nil, fmt.Errorf("error reading circuit pkey: %w", err)
	}
	fmt.Println("pKey loaded, took (s):", time.Since(step))
	step = time.Now()
	// Read and initialize the witness
	cWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		fmt.Println("error initializing witness: ", err)
		return nil, nil, fmt.Errorf("error initializing witness: %w", err)
	}
	if _, err := cWitness.ReadFrom(bytes.NewReader(inputs)); err != nil {
		fmt.Println("error reading witness: ", err)
		return nil, nil, fmt.Errorf("error reading witness: %w", err)
	}
	fmt.Println("witness loaded, took (s):", time.Since(step))

	std.RegisterHints()

	step = time.Now()
	// Generate the proof
	proof, err := groth16.Prove(ccs, provingKey, cWitness)
	if err != nil {
		fmt.Printf("error generating proof: %v\n", err)
		return nil, nil, fmt.Errorf("error generating proof: %w", err)
	}
	fmt.Println("proof generated, took (s):", time.Since(step))
	proofBuff := bytes.Buffer{}
	if _, err := proof.WriteTo(&proofBuff); err != nil {
		return nil, nil, fmt.Errorf("error encoding proof: %w", err)
	}
	// Get public witness part and encode it
	publicWitness, err := cWitness.Public()
	if err != nil {
		return nil, nil, fmt.Errorf("error generating public witness: %w", err)
	}
	publicWitnessBuff := bytes.Buffer{}
	if _, err := publicWitness.WriteTo(&publicWitnessBuff); err != nil {
		return nil, nil, fmt.Errorf("error encoding public witness: %w", err)
	}
	return proofBuff.Bytes(), publicWitnessBuff.Bytes(), nil
}

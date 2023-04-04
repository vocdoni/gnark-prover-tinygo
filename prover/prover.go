package prover

import (
	"bytes"
	"fmt"
	"time"

	"github.com/vocdoni/gnark-crypto-bn254/ecc"
	"github.com/vocdoni/gnark-crypto-bn254/kzg"
	"github.com/vocdoni/gnark-wasm-prover/csbn254"
	"github.com/vocdoni/gnark-wasm-prover/prover"
	"github.com/vocdoni/gnark-wasm-prover/witness"
)

// GenerateProof sets up the circuit with the constrain system and the srs files
// provided and generates the proof for the JSON encoded inputs (witness). It
// returns the verification key, the proof and the public witness, all of this
// outputs will be encoded as JSON. If something fails, it returns an error.
func GenerateProof(bccs, bsrs, bpkey, inputs []byte) ([]byte, []byte, error) {
	step := time.Now()
	// Read and initialize circuit CS
	fmt.Println("loading circuit...")
	ccs := &csbn254.SparseR1CS{}
	if _, err := ccs.ReadFrom(bytes.NewReader(bccs)); err != nil {
		return nil, nil, fmt.Errorf("error reading circuit cs: %w", err)
	}
	fmt.Println("ccs loaded, took (s):", time.Since(step))
	step = time.Now()
	// Read and initialize SSR
	srs := kzg.NewSRS(ecc.BN254)
	if _, err := srs.ReadFrom(bytes.NewReader(bsrs)); err != nil {
		return nil, nil, err
	}
	fmt.Println("srs loaded, took (s):", time.Since(step))
	step = time.Now()
	// Read proving key
	provingKey := &prover.ProvingKey{}
	if _, err := provingKey.ReadFrom(bytes.NewReader(bpkey)); err != nil {
		return nil, nil, fmt.Errorf("error reading circuit pkey: %w", err)
	}
	fmt.Println("pKey loaded, took (s):", time.Since(step))
	step = time.Now()
	// Instance KZG into the proving key
	if err := provingKey.InitKZG(srs); err != nil {
		return nil, nil, fmt.Errorf("error initializating kzg into the pkey: %w", err)
	}
	fmt.Println("kzg initializated into the pKey, took (s):", time.Since(step))
	step = time.Now()
	// Read and initialize the witness
	cWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing witness: %w", err)
	}
	if _, err := cWitness.ReadFrom(bytes.NewReader(inputs)); err != nil {
		return nil, nil, fmt.Errorf("error reading witness: %w", err)
	}
	fmt.Println("witness loaded, took (s):", time.Since(step))
	step = time.Now()

	// Generate the proof
	proof, err := prover.Prove(ccs, provingKey, cWitness)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating proof: %w", err)
	}
	fmt.Println("proof generated, took (s):", time.Since(step))
	proofBuff := bytes.Buffer{}
	if _, err := proof.WriteTo(&proofBuff); err != nil {
		return nil, nil, fmt.Errorf("error encoding proof: %w", err)
	}
	pKeyBuff := bytes.Buffer{}
	if _, err := provingKey.WriteTo(&pKeyBuff); err != nil {
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

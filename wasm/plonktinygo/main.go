package main

import (
	"bytes"
	"fmt"
	"syscall/js"
	"time"

	"github.com/vocdoni/gnark-crypto-bn254/ecc"
	"github.com/vocdoni/gnark-crypto-bn254/kzg"
	"github.com/vocdoni/gnark-wasm-prover/csbn254"
	"github.com/vocdoni/gnark-wasm-prover/prover"
	"github.com/vocdoni/gnark-wasm-prover/witness"
	// This import fixes the issue that raises when a prover tries to generate a proof
	// of a serialized circuit. Check more information here:
	//   - https://github.com/ConsenSys/gnark/issues/600
	//   - https://github.com/phated/gnark-browser/blob/2446c65e89156f1a04163724a89e5dcb7e4c4886/README.md#solution-hint-registration
	// _ "github.com/consensys/gnark/std/math/bits"
)

// GenerateProof sets up the circuit with the constrain system and the srs files
// provided and generates the proof for the JSON encoded inputs (witness). It
// returns the verification key, the proof and the public witness, all of this
// outputs will be encoded as JSON. If something fails, it returns an error.
func GenerateProof(bccs, bsrs, bpkey, inputs []byte) ([]byte, []byte, error) {
	step := time.Now()
	// Read and initialize circuit CS
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

func main() {
	c := make(chan int)
	js.Global().Set("generateProof", js.FuncOf(jsGenerateProof))
	<-c
}

func jsGenerateProof(this js.Value, args []js.Value) interface{} {
	// var bccs, bsrs, witness []byte
	bccs := make([]byte, args[0].Get("length").Int())
	bsrs := make([]byte, args[1].Get("length").Int())
	bpkey := make([]byte, args[2].Get("length").Int())
	bwitness := make([]byte, args[3].Get("length").Int())

	js.CopyBytesToGo(bccs, args[0])
	js.CopyBytesToGo(bsrs, args[1])
	js.CopyBytesToGo(bpkey, args[2])
	js.CopyBytesToGo(bwitness, args[3])

	if _, _, err := GenerateProof(bccs, bsrs, bpkey, bwitness); err != nil {
		return err
	}
	return nil
}

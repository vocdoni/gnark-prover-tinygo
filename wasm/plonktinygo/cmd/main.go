package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/vocdoni/gnark-crypto-bn254/ecc"
	"github.com/vocdoni/gnark-crypto-bn254/kzg"
	"github.com/vocdoni/gnark-wasm-prover/csbn254"
	"github.com/vocdoni/gnark-wasm-prover/hints"
	"github.com/vocdoni/gnark-wasm-prover/prover"
	"github.com/vocdoni/gnark-wasm-prover/witness"
	// This import fixes the issue that raises when a prover tries to generate a proof
	// of a serialized circuit. Check more information here:
	//   - https://github.com/ConsenSys/gnark/issues/600
	//   - https://github.com/phated/gnark-browser/blob/2446c65e89156f1a04163724a89e5dcb7e4c4886/README.md#solution-hint-registration
	// "github.com/consensys/gnark/std/math/bits"
)

// GenerateProof sets up the circuit with the constrain system and the srs files
// provided and generates the proof for the JSON encoded inputs (witness). It
// returns the verification key, the proof and the public witness, all of this
// outputs will be encoded as JSON. If something fails, it returns an error.
//
//export generateProof
func GenerateProof(bccs, bsrs, bpkey, inputs io.Reader) ([]byte, []byte, error) {
	step := time.Now()
	// Read and initialize circuit CS
	ccs := &csbn254.SparseR1CS{}
	if _, err := ccs.ReadFrom(bccs); err != nil {
		return nil, nil, fmt.Errorf("error reading circuit cs: %w", err)
	}
	fmt.Println("ccs loaded, took (s):", time.Since(step))
	//fmt.Printf("\n\nvar ccs = %#v\n\n", ccs)
	step = time.Now()
	// Read and initialize SSR
	srs := kzg.NewSRS(ecc.BN254)
	if _, err := srs.ReadFrom(bsrs); err != nil {
		return nil, nil, err
	}
	fmt.Println("srs loaded, took (s):", time.Since(step))
	step = time.Now()
	// Read proving key
	provingKey := &prover.ProvingKey{}
	if _, err := provingKey.ReadFrom(bpkey); err != nil {
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
	if _, err := cWitness.ReadFrom(inputs); err != nil {
		return nil, nil, fmt.Errorf("error reading witness: %w", err)
	}
	fmt.Println("witness loaded, took (s):", time.Since(step))
	step = time.Now()

	// Register hints for the circuit
	hints.RegisterHints()

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
	fdcircuit := flag.String("circuit", "", "circuit file")
	fdsrs := flag.String("srs", "", "srs file")
	fdpkey := flag.String("pkey", "", "proving key file")
	fdwitness := flag.String("witness", "", "witness file")

	flag.Parse()

	// Read the files into byte slices and call the generateProof function
	fmt.Println("reading circuit file: ", *fdcircuit)
	fccs, err := os.Open(*fdcircuit)
	if err != nil {
		panic(err)
	}
	defer fccs.Close()

	fmt.Println("reading srs file: ", *fdsrs)
	fsrs, err := os.Open(*fdsrs)
	if err != nil {
		panic(err)
	}

	defer fsrs.Close()
	fmt.Println("reading proving key file: ", *fdpkey)
	fpkey, err := os.Open(*fdpkey)
	if err != nil {
		panic(err)
	}
	defer fpkey.Close()

	fmt.Println("reading witness file: ", *fdwitness)
	fwitness, err := os.Open(*fdwitness)
	if err != nil {
		panic(err)
	}
	defer fwitness.Close()

	fmt.Println("calling generateProof function")
	proof, publicWitness, err := GenerateProof(fccs, fsrs, fpkey, fwitness)
	if err != nil {
		panic(err)
	}
	fmt.Printf("proof: %x\n", proof)
	fmt.Printf("public witness: %x\n", publicWitness)
}

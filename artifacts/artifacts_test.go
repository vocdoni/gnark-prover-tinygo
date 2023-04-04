package artifacts

import (
	"bytes"
	"os"
	"testing"

	"gnark-prover-tinygo/prover"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	qt "github.com/frankban/quicktest"
)

func TestArtifacts(t *testing.T) {
	c := qt.New(t)
	// load ccs, srs, pkey, vkey and witness artifacts
	bccs, err := os.ReadFile("./zkcensus.ccs")
	c.Assert(err, qt.IsNil)
	// dismiss error opening srs (not required with groth16)
	bsrs, _ := os.ReadFile("./zkcensus.srs")
	bpkey, err := os.ReadFile("./zkcensus.pkey")
	c.Assert(err, qt.IsNil)
	bvkey, err := os.ReadFile("./zkcensus.vkey")
	c.Assert(err, qt.IsNil)
	bwitness, err := os.ReadFile("./zkcensus.witness")
	c.Assert(err, qt.IsNil)
	// generate proof with plonk and verify it
	bproof, bpubwitness, err := prover.GenerateProof(bccs, bsrs, bpkey, bwitness)
	c.Assert(err, qt.IsNil)
	// parse the verifiying key
	verifiyingKey := plonk.NewVerifyingKey(ecc.BN254)
	_, err = verifiyingKey.ReadFrom(bytes.NewBuffer(bvkey))
	c.Assert(err, qt.IsNil)
	// Read and initialize SSR
	srs := kzg.NewSRS(ecc.BN254)
	_, err = srs.ReadFrom(bytes.NewReader(bsrs))
	c.Assert(err, qt.IsNil)
	err = verifiyingKey.InitKZG(srs)
	c.Assert(err, qt.IsNil)
	// parse the proof
	proof := plonk.NewProof(ecc.BN254)
	_, err = proof.ReadFrom(bytes.NewBuffer(bproof))
	c.Assert(err, qt.IsNil)
	// parse the public witness
	pubWitness, err := witness.New(ecc.BN254.ScalarField())
	c.Assert(err, qt.IsNil)
	_, err = pubWitness.ReadFrom(bytes.NewReader(bpubwitness))
	c.Assert(err, qt.IsNil)
	// assert the result of the verification process
	err = plonk.Verify(proof, verifiyingKey, pubWitness)
	c.Assert(err, qt.IsNil)
}

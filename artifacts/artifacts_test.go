package artifacts

import (
	"gnark-prover-tinygo/internal/circuit/groth16"
	"gnark-prover-tinygo/internal/circuit/plonk"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestArtifacts(t *testing.T) {
	c := qt.New(t)

	ccs, err := os.ReadFile("./zkcensus.ccs")
	c.Assert(err, qt.IsNil)

	srs, err := os.ReadFile("./zkcensus.srs")
	c.Assert(err, qt.IsNil)

	witness, err := os.ReadFile("witness")
	c.Assert(err, qt.IsNil)

	// Try Plonk artifacts
	vk, proof, pubWitness, err := plonk.GenerateProof(ccs, srs, witness)
	if err == nil {
		err = plonk.VerifyProof(srs, vk, proof, pubWitness)
		c.Assert(err, qt.IsNil)
		return
	}

	// Try Groth16 artifacts
	vk, proof, pubWitness, err = groth16.GenerateProof(ccs, srs, witness)
	c.Assert(err, qt.IsNil)

	err = groth16.VerifyProof(srs, vk, proof, pubWitness)
	c.Assert(err, qt.IsNil)
}

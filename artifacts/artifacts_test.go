package artifacts

import (
	"gnark-prover-tinygo/internal/circuit"
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

	vk, proof, pubWitness, err := circuit.GenerateProof(ccs, srs, witness)
	c.Assert(err, qt.IsNil)

	err = circuit.VerifyProof(srs, vk, proof, pubWitness)
	c.Assert(err, qt.IsNil)
}

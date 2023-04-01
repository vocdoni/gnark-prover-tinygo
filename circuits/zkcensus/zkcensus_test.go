package zkcensus

import (
	"flag"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/test"
)

var nLevels = flag.Int("nLevels", 160, "number of levels of the arbo tree")
var nKeys = flag.Int("nKyes", 200, "number of keys to add to the arbo tree")

func TestZkCensusCircuit(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit ZkCensusCircuit

	success, _ := ZkCensusInputs(*nLevels, *nKeys)
	assert.SolvingSucceeded(&circuit, &success, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
	assert.SolvingSucceeded(&circuit, &success, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))
}

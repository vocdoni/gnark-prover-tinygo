package smt

import (
	"gnark-test/internal/zkaddress"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	qt "github.com/frankban/quicktest"
	"go.vocdoni.io/dvote/db"
	"go.vocdoni.io/dvote/db/pebbledb"
	"go.vocdoni.io/dvote/tree/arbo"
	"go.vocdoni.io/dvote/util"
)

type testVerifierCircuit struct {
	Root     frontend.Variable
	Key      frontend.Variable
	Value    frontend.Variable
	Siblings [160]frontend.Variable
}

func (circuit *testVerifierCircuit) Define(api frontend.API) error {
	return Verifier(api, circuit.Root, circuit.Key, circuit.Value, circuit.Siblings[:])
}

func successInputs(t *testing.T) testVerifierCircuit {
	c := qt.New(t)

	database, err := pebbledb.New(db.Options{Path: t.TempDir()})
	c.Assert(err, qt.IsNil)
	arboTree, err := arbo.NewTree(arbo.Config{
		Database:     database,
		MaxLevels:    160,
		HashFunction: arbo.HashFunctionPoseidon,
	})
	c.Assert(err, qt.IsNil)

	factoryWeight := big.NewInt(10)
	candidate, err := zkaddress.FromBytes(util.RandomBytes(32))
	c.Assert(err, qt.IsNil)
	err = arboTree.Add(candidate.ArboBytes(), arbo.BigIntToBytes(arbo.HashFunctionPoseidon.Len(), factoryWeight))
	c.Assert(err, qt.IsNil)

	for i := 0; i < 100; i++ {
		k, err := zkaddress.FromBytes(util.RandomBytes(32))
		c.Assert(err, qt.IsNil)
		err = arboTree.Add(k.ArboBytes(), arbo.BigIntToBytes(arbo.HashFunctionPoseidon.Len(), factoryWeight))
		c.Assert(err, qt.IsNil)
	}

	key, value, pSiblings, exist, err := arboTree.GenProof(candidate.ArboBytes())
	c.Assert(err, qt.IsNil)
	c.Assert(exist, qt.IsTrue)
	c.Assert(key, qt.ContentEquals, candidate.ArboBytes())
	c.Assert(value, qt.ContentEquals, arbo.BigIntToBytes(arbo.HashFunctionPoseidon.Len(), factoryWeight))

	uSiblings, err := arbo.UnpackSiblings(arbo.HashFunctionPoseidon, pSiblings)
	c.Assert(err, qt.IsNil)

	siblings := [160]frontend.Variable{}
	for i := 0; i < 160; i++ {
		if i < len(uSiblings) {
			siblings[i] = arbo.BytesToBigInt(uSiblings[i])
		} else {
			siblings[i] = big.NewInt(0)
		}
	}

	root, err := arboTree.Root()
	c.Assert(err, qt.IsNil)
	return testVerifierCircuit{
		Root:     arbo.BytesToBigInt(root),
		Key:      candidate.Scalar,
		Value:    factoryWeight,
		Siblings: siblings,
	}
}

func TestVerifier(t *testing.T) {
	assert := test.NewAssert(t)

	var circuit testVerifierCircuit
	inputs := successInputs(t)
	assert.SolvingSucceeded(&circuit, &inputs, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
}

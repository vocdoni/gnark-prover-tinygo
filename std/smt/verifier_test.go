package smt

import (
	"gnark-prover-tinygo/internal/arbo"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	qt "github.com/frankban/quicktest"
	"go.vocdoni.io/dvote/db"
	"go.vocdoni.io/dvote/db/pebbledb"
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
	// set number of levels and max key and value length
	nLevels := 160
	kLen := nLevels / 8
	// create database
	database, err := pebbledb.New(db.Options{Path: t.TempDir()})
	c.Assert(err, qt.IsNil)
	// instance a new arbo tree
	tree, err := arbo.NewTree(arbo.Config{
		Database:     database,
		MaxLevels:    nLevels,
		HashFunction: arbo.HashFunctionMiMC,
	})
	c.Assert(err, qt.IsNil)
	c.Assert(err, qt.IsNil)
	// add 100 more random keys to the arbo tree with the default value
	for i := 0; i < 100; i++ {
		k := arbo.BigIntToBytes(kLen, big.NewInt(int64(i)))
		v := arbo.BigIntToBytes(kLen, big.NewInt(int64(i*2)))
		err = tree.Add(k, v)
		c.Assert(err, qt.IsNil)
	}
	// get and encode the merkle root
	root, err := tree.Root()
	c.Assert(err, qt.IsNil)
	// instance proof candidate key and default value
	k := arbo.BigIntToBytes(kLen, new(big.Int).SetInt64(7))
	v := arbo.BigIntToBytes(kLen, big.NewInt(14))
	// generate the proof for the candidate key
	kAux, vAux, pSiblings, exist, err := tree.GenProof(k)
	c.Assert(err, qt.IsNil)
	c.Assert(exist, qt.IsTrue)
	c.Assert(kAux, qt.DeepEquals, k)
	c.Assert(vAux, qt.DeepEquals, v)
	// check the generated proof
	valid, err := arbo.CheckProof(tree.HashFunction(), k, v, root, pSiblings)
	c.Assert(err, qt.IsNil)
	c.Assert(valid, qt.IsTrue)
	// unpack the proof siblings
	uSiblings, err := arbo.UnpackSiblings(tree.HashFunction(), pSiblings)
	c.Assert(err, qt.IsNil)
	// encode the siblings into a array of circuit inputs
	siblings := [160]frontend.Variable{}
	for i := 0; i < 160; i++ {
		if i < len(uSiblings) {
			siblings[i] = arbo.BytesToBigInt(uSiblings[i])
		} else {
			siblings[i] = big.NewInt(0)
		}
	}
	return testVerifierCircuit{
		Root:     arbo.BytesToBigInt(root),
		Key:      arbo.BytesToBigInt(k),
		Value:    arbo.BytesToBigInt(v),
		Siblings: siblings,
	}
}

func TestVerifier(t *testing.T) {
	assert := test.NewAssert(t)

	var circuit testVerifierCircuit
	inputs := successInputs(t)
	assert.SolvingSucceeded(&circuit, &inputs, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
}

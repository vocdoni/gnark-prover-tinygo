package zkcensus

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"gnark-test/internal/zkaddress"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	qt "github.com/frankban/quicktest"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"go.vocdoni.io/dvote/db"
	"go.vocdoni.io/dvote/db/pebbledb"
	"go.vocdoni.io/dvote/tree/arbo"
	"go.vocdoni.io/dvote/util"
)

var nLevels = flag.Int("nLevels", 160, "number of levels of the arbo tree")
var nKeys = flag.Int("nKyes", 200, "number of keys to add to the arbo tree")

func emptyInput() ZkCensusCircuit {
	return ZkCensusCircuit{
		ElectionId:     [2]frontend.Variable{0, 0},
		CensusRoot:     frontend.Variable(0),
		Nullifier:      frontend.Variable(0),
		FactoryWeight:  frontend.Variable(0),
		VoteHash:       [2]frontend.Variable{0, 0},
		CensusSiblings: [160]frontend.Variable{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		PrivateKey:     frontend.Variable(0),
		VotingWeight:   frontend.Variable(0),
	}
}

func BytesToArbo(input []byte) [2]*big.Int {
	hash := sha256.Sum256(input)
	return [2]*big.Int{
		new(big.Int).SetBytes(arbo.SwapEndianness(hash[:16])),
		new(big.Int).SetBytes(arbo.SwapEndianness(hash[16:])),
	}
}

func correctInputs(t *testing.T) ZkCensusCircuit {
	c := qt.New(t)
	database, err := pebbledb.New(db.Options{Path: t.TempDir()})
	c.Assert(err, qt.IsNil)

	arboTree, err := arbo.NewTree(arbo.Config{
		Database:     database,
		MaxLevels:    *nLevels,
		HashFunction: arbo.HashFunctionPoseidon,
	})
	c.Assert(err, qt.IsNil)

	mockValue := big.NewInt(10)
	candidate, err := zkaddress.FromBytes([]byte("1b505cdafb4b1150b1a740633af41e5e1f19a5c4"))

	fmt.Println(candidate.Scalar)
	fmt.Println(candidate.Public)
	c.Assert(err, qt.IsNil)

	err = arboTree.Add(candidate.Scalar.Bytes(), mockValue.Bytes())
	c.Assert(err, qt.IsNil)

	for i := 1; i < *nKeys; i++ {
		k, err := zkaddress.FromBytes([]byte(util.RandomHex(32)))
		c.Assert(err, qt.IsNil)
		err = arboTree.Add(k.Scalar.Bytes(), mockValue.Bytes())
		c.Assert(err, qt.IsNil)
	}

	key, value, pSiblings, exist, err := arboTree.GenProof(candidate.Scalar.Bytes())
	c.Assert(err, qt.IsNil)
	c.Assert(exist, qt.IsTrue)
	c.Assert(key, qt.ContentEquals, candidate.Scalar.Bytes())
	c.Assert(value, qt.ContentEquals, mockValue.Bytes())

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
	censusRoot := arbo.BytesToBigInt(root)

	electionId := BytesToArbo([]byte(util.RandomHex(32)))
	nullifier, err := poseidon.Hash([]*big.Int{candidate.Private, electionId[0], electionId[1]})
	c.Assert(err, qt.IsNil)

	voteHash := BytesToArbo(mockValue.Bytes())
	c.Assert(err, qt.IsNil)
	return ZkCensusCircuit{
		ElectionId:     [2]frontend.Variable{electionId[0], electionId[1]},
		CensusRoot:     censusRoot,
		Nullifier:      nullifier,
		FactoryWeight:  mockValue,
		VoteHash:       [2]frontend.Variable{voteHash[0], voteHash[1]},
		CensusSiblings: siblings,
		PrivateKey:     candidate.Private,
		VotingWeight:   big.NewInt(5),
	}
}

func TestZkCensusCircuit(t *testing.T) {
	assert := test.NewAssert(t)

	var circuit ZkCensusCircuit

	fail := emptyInput()
	assert.SolvingFailed(&circuit, &fail, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))

	success := correctInputs(t)
	assert.SolvingSucceeded(&circuit, &success, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
}

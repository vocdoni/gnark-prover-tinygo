package smtverifier

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"testing"

	"go.vocdoni.io/dvote/db"
	"go.vocdoni.io/dvote/db/pebbledb"
	"go.vocdoni.io/dvote/tree/arbo"
	"go.vocdoni.io/dvote/util"

	qt "github.com/frankban/quicktest"
)

var nLevels = flag.Int("nLevels", 160, "number of levels of the arbo tree")
var nKeys = flag.Int("nKyes", 200, "number of keys to add to the arbo tree")

func randKey() *big.Int {
	return new(big.Int).SetInt64(int64(util.RandomInt(100000, 10000000)))
}

type ArboProof struct {
	Root     *big.Int   `json:"Root"`
	Key      *big.Int   `json:"Key"`
	Value    *big.Int   `json:"Value"`
	Siblings []*big.Int `json:"Siblings"`
}

func TestSMTVerifierInputs(t *testing.T) {
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
	candidate, _ := new(big.Int).SetString("9855069924893712724378342796031175650258250494", 10)
	keys := [][]byte{candidate.Bytes()}
	err = arboTree.Add(candidate.Bytes(), mockValue.Bytes())
	c.Assert(err, qt.IsNil)

	for i := 1; i < *nKeys; i++ {
		k := randKey().Bytes()
		keys = append(keys, k)
		err = arboTree.Add(k, mockValue.Bytes())
		c.Assert(err, qt.IsNil)
	}

	key, value, pSiblings, exist, err := arboTree.GenProof(candidate.Bytes())
	c.Assert(err, qt.IsNil)
	c.Assert(exist, qt.IsTrue)
	c.Assert(key, qt.ContentEquals, candidate.Bytes())
	c.Assert(value, qt.ContentEquals, mockValue.Bytes())

	uSiblings, err := arbo.UnpackSiblings(arbo.HashFunctionPoseidon, pSiblings)
	c.Assert(err, qt.IsNil)

	siblings := []*big.Int{}
	for i := 0; i < len(uSiblings); i++ {
		siblings = append(siblings, arbo.BytesToBigInt(uSiblings[i]))
	}

	root, err := arboTree.Root()
	c.Assert(err, qt.IsNil)

	ok, err := arbo.CheckProof(arbo.HashFunctionPoseidon, candidate.Bytes(), mockValue.Bytes(), root, pSiblings)
	c.Assert(err, qt.IsNil)
	c.Assert(ok, qt.IsTrue)

	result := ArboProof{
		Key:      candidate,
		Root:     arbo.BytesToBigInt(root),
		Value:    mockValue,
		Siblings: siblings,
	}
	pResult, err := json.MarshalIndent(result, "", "    ")
	c.Assert(err, qt.IsNil)
	fmt.Println(string(pResult))
}

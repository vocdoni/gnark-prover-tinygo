package zkcensus

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"gnark-prover-tinygo/internal/zkaddress"
	"math/big"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
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

func correctInputs() (ZkCensusCircuit, error) {
	dbTemp, err := os.MkdirTemp("", "db")
	if err != nil {
		return ZkCensusCircuit{}, err
	}
	database, err := pebbledb.New(db.Options{Path: dbTemp})
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	arboTree, err := arbo.NewTree(arbo.Config{
		Database:     database,
		MaxLevels:    *nLevels,
		HashFunction: arbo.HashFunctionPoseidon,
	})
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	factoryWeight := big.NewInt(10)
	encFactoryWeight := arbo.BigIntToBytes(arbo.HashFunctionPoseidon.Len(), factoryWeight)
	candidate, err := zkaddress.FromBytes([]byte("1b505cdafb4b1150b1a740633af41e5e1f19a5c4"))
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	err = arboTree.Add(candidate.ArboBytes(), encFactoryWeight)
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	for i := 1; i < *nKeys; i++ {
		k, err := zkaddress.FromBytes(util.RandomBytes(32))
		if err != nil {
			return ZkCensusCircuit{}, err
		}

		err = arboTree.Add(k.ArboBytes(), encFactoryWeight)
		if err != nil {
			return ZkCensusCircuit{}, err
		}
	}

	_, _, pSiblings, exist, err := arboTree.GenProof(candidate.ArboBytes())
	if err != nil {
		return ZkCensusCircuit{}, err
	} else if !exist {
		return ZkCensusCircuit{}, fmt.Errorf("key does not exists")
	}

	uSiblings, err := arbo.UnpackSiblings(arbo.HashFunctionPoseidon, pSiblings)
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	siblings := [160]frontend.Variable{}
	for i := 0; i < 160; i++ {
		if i < len(uSiblings) {
			siblings[i] = arbo.BytesToBigInt(uSiblings[i])
		} else {
			siblings[i] = big.NewInt(0)
		}
	}

	root, err := arboTree.Root()
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	electionId := BytesToArbo(util.RandomBytes(32))
	nullifier, err := poseidon.Hash([]*big.Int{candidate.Private, electionId[0], electionId[1]})
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	voteHash := BytesToArbo(factoryWeight.Bytes())
	return ZkCensusCircuit{
		ElectionId:     [2]frontend.Variable{electionId[0], electionId[1]},
		CensusRoot:     arbo.BytesToBigInt(root),
		Nullifier:      nullifier,
		FactoryWeight:  factoryWeight,
		VoteHash:       [2]frontend.Variable{voteHash[0], voteHash[1]},
		CensusSiblings: siblings,
		PrivateKey:     candidate.Private,
		VotingWeight:   big.NewInt(5),
	}, nil
}

func SerializeWitness() error {
	success, err := correctInputs()
	if err != nil {
		return err
	}
	witness, _ := frontend.NewWitness(&success, ecc.BN254.ScalarField())
	f, err := os.Create("./witness")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = witness.WriteTo(f)
	return err
}

func TestZkCensusCircuit(t *testing.T) {
	assert := test.NewAssert(t)

	var circuit ZkCensusCircuit

	fail := emptyInput()
	assert.SolvingFailed(&circuit, &fail, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))

	success, _ := correctInputs()
	assert.SolvingSucceeded(&circuit, &success, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
}

package zkcensus

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"gnark-prover-tinygo/internal/arbo"
	"gnark-prover-tinygo/internal/zkaddress"
	"math/big"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"go.vocdoni.io/dvote/db"
	"go.vocdoni.io/dvote/db/pebbledb"

	"go.vocdoni.io/dvote/util"
)

var nLevels = flag.Int("nLevels", 160, "number of levels of the arbo tree")
var nKeys = flag.Int("nKyes", 200, "number of keys to add to the arbo tree")

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
		HashFunction: arbo.HashFunctionMiMC,
	})
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	factoryWeight := big.NewInt(10)
	encFactoryWeight := arbo.BigIntToBytes(arbo.HashFunctionMiMC.Len(), factoryWeight)
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

	root, err := arboTree.Root()
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	if valid, err := arbo.CheckProof(arboTree.HashFunction(), candidate.ArboBytes(), encFactoryWeight, root, pSiblings); err != nil {
		return ZkCensusCircuit{}, err
	} else if !valid {
		return ZkCensusCircuit{}, fmt.Errorf("proof not valid")
	}

	uSiblings, err := arbo.UnpackSiblings(arbo.HashFunctionMiMC, pSiblings)
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

	electionId := BytesToArbo(util.RandomBytes(32))
	var hash arbo.HashMiMC

	bNullifier, err := hash.Hash(
		arbo.BigIntToBytes(arboTree.HashFunction().Len(), candidate.Private),
		arbo.BigIntToBytes(arboTree.HashFunction().Len(), electionId[0]),
		arbo.BigIntToBytes(arboTree.HashFunction().Len(), electionId[1]),
	)
	if err != nil {
		return ZkCensusCircuit{}, err
	}

	voteHash := BytesToArbo(factoryWeight.Bytes())
	return ZkCensusCircuit{
		ElectionId:     [2]frontend.Variable{electionId[0], electionId[1]},
		CensusRoot:     arbo.BytesToBigInt(root),
		Nullifier:      arbo.BytesToBigInt(bNullifier),
		FactoryWeight:  factoryWeight,
		VoteHash:       [2]frontend.Variable{voteHash[0], voteHash[1]},
		CensusSiblings: siblings,
		PrivateKey:     candidate.Private,
		VotingWeight:   big.NewInt(5),
	}, nil
}

func SerializeWitness(input ZkCensusCircuit) error {
	witness, _ := frontend.NewWitness(&input, ecc.BN254.ScalarField())
	f, err := os.Create("../artifacts/witness")
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

	success, _ := correctInputs()
	SerializeWitness(success)
	assert.SolvingSucceeded(&circuit, &success, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
	assert.SolvingSucceeded(&circuit, &success, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))
}

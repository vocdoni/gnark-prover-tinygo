package zkcensus

import (
	"gnark-test/std/hash/poseidon"
	"gnark-test/std/smt"
	"gnark-test/std/zkaddress"

	"github.com/consensys/gnark/frontend"
)

type ZkCensusCircuit struct {
	// Public inputs
	ElectionId    [2]frontend.Variable
	CensusRoot    frontend.Variable
	Nullifier     frontend.Variable
	FactoryWeight frontend.Variable
	VoteHash      [2]frontend.Variable
	// Private inputs
	CensusSiblings [160]frontend.Variable
	PrivateKey     frontend.Variable
	VotingWeight   frontend.Variable
}

func (circuit *ZkCensusCircuit) Define(api frontend.API) error {
	// votingWeight represents the weight that the user wants to use to perform
	// a vote and must be lower than factoryWeight
	api.AssertIsLessOrEqual(circuit.VotingWeight, circuit.FactoryWeight)
	// voteHash is not operated inside the circuit, assuming that in
	// Circom an input that is not used will be included in the constraints
	// system and in the witness
	api.AssertIsDifferent(circuit.VoteHash[0], 0)
	api.AssertIsDifferent(circuit.VoteHash[1], 0)
	// calculate the zkaddress from the private key
	zkAddr, err := zkaddress.FromPrivate(api, circuit.PrivateKey)
	if err != nil {
		return err
	}

	api.Println(zkAddr.Scalar)
	api.Println(zkAddr.Public)
	// check the Merkletree with census root, siblings, zkAddress and factory
	// weight
	if err := smt.Verifier(api, circuit.CensusRoot, zkAddr.Scalar, circuit.FactoryWeight, circuit.CensusSiblings[:]); err != nil {
		return err
	}
	// check nullifier (electionID + privateKey)
	computedNullifier := poseidon.Hash(api, circuit.PrivateKey, circuit.ElectionId[0], circuit.ElectionId[1])
	api.AssertIsEqual(circuit.Nullifier, computedNullifier)

	return nil
}

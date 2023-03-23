package zkcensus

import (
	"gnark-prover-tinygo/std/hash/poseidon"
	"gnark-prover-tinygo/std/smt"
	"gnark-prover-tinygo/std/zkaddress"

	"github.com/consensys/gnark/frontend"
)

type ZkCensusCircuit struct {
	// Public inputs
	ElectionId    [2]frontend.Variable `gnark:",public"`
	CensusRoot    frontend.Variable    `gnark:",public"`
	Nullifier     frontend.Variable    `gnark:",public"`
	FactoryWeight frontend.Variable    `gnark:",public"`
	VoteHash      [2]frontend.Variable `gnark:",public"`
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

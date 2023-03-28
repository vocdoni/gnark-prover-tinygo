/*
ZkCensus is the Vocdoni ZKSnark circuit to proof that a given voter is part of a
created election census (in a Arbo Merkle Tree). The circuit checks:

	-the prover is the owner of the private key
	-keyHash (hash of the user's public key) belongs to the census
		-the public key is generated based on the provided private key
		-the public key is inside a hash, which is inside the Merkletree with
		the CensusRoot and siblings (key=keyHash, value=factoryWeight)
	-H(private key, processID) == nullifier
		-to avoid proof reusability
	-factoryWeight is the weight assigned by default to the owner of the private
	key andincluded as merkle tree leaf value.
	-votingWeight is the weight desired to vote by the owner of the private key
	and must be less than or equal to the factoryWeightht.
*/
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

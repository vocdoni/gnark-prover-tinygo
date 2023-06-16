package prover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/std"
	gnarkparser "github.com/vocdoni/go-snark/parsers/gnark"
	"github.com/vocdoni/go-snark/prover"

	// This import fixes the issue that raises when a prover tries to generate a proof
	// of a serialized circuit. Check more information here:
	//   - https://github.com/ConsenSys/gnark/issues/600
	//   - https://github.com/phated/gnark-browser/blob/2446c65e89156f1a04163724a89e5dcb7e4c4886/README.md#solution-hint-registration
	_ "github.com/consensys/gnark/std/math/bits"
)

/*type ZkCensusCircuit struct {
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

var circuitInputs = ZkCensusCircuit{
	ElectionId:    [2]frontend.Variable{"133469402295578115674590583255628297017", "238408956434793107449390126332638610245"},
	CensusRoot:    "18791043466379219967070787962573051007292340205045330061714942310370017591784",
	Nullifier:     "11497068997553757467335284958718269729876301997209500898124318310115819586257",
	FactoryWeight: 10,
	VoteHash:      [2]frontend.Variable{"242108076058607163538102198631955675649", "142667662805314151155817304537028292174"},
	PrivateKey:    "12007696602022466067210558438468234995085206818257350359618361229442198701667",
	VotingWeight:  5,
	CensusSiblings: [160]frontend.Variable{
		"3750341641233045923158109099254819542097988317102271429655215925695405823874",
		"18945719540214410315511518748203941019347214565125813490891947393805931340459",
		"17433150599005607441584191131943920757473796641236161857902215334239628633319",
		"12338881141001312372400639848762285844832728978091440493928026493220985914320",
		"19730049984892730507117383126083647386940053917094166303449215211583530399428",
		"5329573256950550240549165409825997834544442018148185937946334839460703646733",
		"1736424198595143520428995953295401501339874905590394281091009008029487076353",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0",
		"0", "0", "0", "0", "0", "0", "0", "0", "0"},
}
*/

func GenerateProofGroth16GoSnark(pkBin, witBin []byte) ([]byte, []byte, error) {
	pk, err := gnarkparser.TransformProvingKey(bytes.NewReader(pkBin))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing proving key: %w", err)
	}
	wit, err := gnarkparser.TransformWitness(bytes.NewReader(witBin))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing witness: %w", err)
	}

	/*fmt.Println("creating witness")
	witness, err := frontend.NewWitness(&circuitInputs, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, fmt.Errorf("error creating witness: %w", err)
	}
	wbuf := bytes.Buffer{}
	wSize, err := witness.WriteTo(&wbuf)
	if err != nil {
		return nil, nil, fmt.Errorf("error writing witness: %w", err)
	}
	fmt.Printf("witness size: %d\n", wSize)
	wit, err := gnarkparser.TransformWitness(bytes.NewReader(wbuf.Bytes()))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing witness: %w", err)
	}
	*/
	start := time.Now()
	proof, pubInputs, err := prover.GenerateProof(pk, *wit)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating proof: %w", err)
	}
	proofData, err := proof.MarshalJSON()
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling proof: %w", err)
	}
	pubInputsSlice := []string{}
	for _, p := range pubInputs {
		pubInputsSlice = append(pubInputsSlice, p.String())
	}
	pubInputsData, err := json.Marshal(pubInputsSlice)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling public inputs: %w", err)
	}

	fmt.Println("proof generation took (s):", time.Since(start))
	fmt.Printf("proof: %s\n", proofData)
	fmt.Printf("pubInputs: %s\n", pubInputsData)

	return proofData, pubInputsData, nil
}

// GenerateProofGroth16 sets up the circuit with the constrain system and the srs files
// provided and generates the proof for the JSON encoded inputs (witness). It
// returns the verification key, the proof and the public witness, all of this
// outputs will be encoded as JSON. If something fails, it returns an error.
func GenerateProofGroth16(bccs, bpkey, inputs []byte) ([]byte, []byte, error) {
	step := time.Now()
	// Read and initialize circuit CS
	ccs := groth16.NewCS(ecc.BN254)
	if _, err := ccs.ReadFrom(bytes.NewReader(bccs)); err != nil {
		fmt.Println("error reading circuit cs: ", err)
		return nil, nil, fmt.Errorf("error reading circuit cs: %w", err)
	}
	fmt.Println("ccs loaded, took (s):", time.Since(step))
	step = time.Now()
	// Read proving key
	provingKey := groth16.NewProvingKey(ecc.BN254)
	if _, err := provingKey.UnsafeReadFrom(bytes.NewReader(bpkey)); err != nil {
		fmt.Println("error reading circuit pkey: ", err)
		return nil, nil, fmt.Errorf("error reading circuit pkey: %w", err)
	}
	fmt.Println("pKey loaded, took (s):", time.Since(step))
	step = time.Now()
	// Read and initialize the witness
	cWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		fmt.Println("error initializing witness: ", err)
		return nil, nil, fmt.Errorf("error initializing witness: %w", err)
	}
	if _, err := cWitness.ReadFrom(bytes.NewReader(inputs)); err != nil {
		fmt.Println("error reading witness: ", err)
		return nil, nil, fmt.Errorf("error reading witness: %w", err)
	}
	fmt.Println("witness loaded, took (s):", time.Since(step))

	std.RegisterHints()

	step = time.Now()
	// Generate the proof
	proof, err := groth16.Prove(ccs, provingKey, cWitness)
	if err != nil {
		fmt.Printf("error generating proof: %v\n", err)
		return nil, nil, fmt.Errorf("error generating proof: %w", err)
	}
	fmt.Println("proof generated, took (s):", time.Since(step))
	proofBuff := bytes.Buffer{}
	if _, err := proof.WriteTo(&proofBuff); err != nil {
		return nil, nil, fmt.Errorf("error encoding proof: %w", err)
	}
	// Get public witness part and encode it
	publicWitness, err := cWitness.Public()
	if err != nil {
		return nil, nil, fmt.Errorf("error generating public witness: %w", err)
	}
	publicWitnessBuff := bytes.Buffer{}
	if _, err := publicWitness.WriteTo(&publicWitnessBuff); err != nil {
		return nil, nil, fmt.Errorf("error encoding public witness: %w", err)
	}
	return proofBuff.Bytes(), publicWitnessBuff.Bytes(), nil
}

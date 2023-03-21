package smt

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
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

func TestVerifier(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit, assignment testVerifierCircuit

	assignment.Root, _ = new(big.Int).SetString("19386423483289217751546547558510578991764514630936230601846549184344089393993", 10)
	assignment.Key, _ = new(big.Int).SetString("416698148726288012817560452240299996958781324757", 10)
	assignment.Value, _ = new(big.Int).SetString("10", 10)

	encSiblings := []string{
		"10538390628779904175123812771770207693268525215246857421995684443828602555706",
		"16604294849135879097828914426213854703737930909791927988940580247725520584905",
		"3391303419572454436552609940337733365298689239216434310424694230130019699858",
		"15626866731086404336676273880574542987417791193566692617410850039682622927975",
		"994142336116276929888831522189142863397704914551242792615621701543904253302",
		"18158941423566262345714473584027873922823703805582428224021578191278099022800",
		"12994893564173120417559034184534249165082426830072306729018045728435876153884",
		"0",
		"12285948963496995327859002909484453104525104463644625648271148480706634647025"}

	assignment.Siblings = [160]frontend.Variable{}
	for i := 0; i < 160; i++ {
		if i < len(encSiblings) {
			assignment.Siblings[i], _ = new(big.Int).SetString(encSiblings[i], 10)
		} else {
			assignment.Siblings[i] = big.NewInt(0)
		}
	}

	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
}

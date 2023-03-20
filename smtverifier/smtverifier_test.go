package smtverifier

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

type testSMTVerifierCircuit struct {
	Root     frontend.Variable
	Key      frontend.Variable
	Value    frontend.Variable
	Siblings [160]frontend.Variable
}

func (circuit *testSMTVerifierCircuit) Define(api frontend.API) error {
	return SMTVerifier(api,
		1,             // enabled
		circuit.Root,  // root
		0,             // oldKey
		0,             // oldValue
		0,             // isOld0
		circuit.Key,   // key
		circuit.Value, // value
		0,             // fnc
		circuit.Siblings[:])
}

func TestSMTVerifier(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit, assignment testSMTVerifierCircuit

	assignment.Root, _ = new(big.Int).SetString("10881694218522057830851846252572493962591578683124701053886591513468217142401", 10)
	assignment.Key, _ = new(big.Int).SetString("9855069924893712724378342796031175650258250494", 10)
	assignment.Value, _ = new(big.Int).SetString("10", 10)

	encSiblings := []string{
		"13270665487936970877834115610350183460918350319317656520964516901798708837389",
		"15754424946272227113852799828604833292482456962450981445711363050808177117942",
		"6849698160061539396610595734152752655189502140332503219979718656598455018898",
		"1498443836521874323689010039765882270347698540745288184716404943008071063256",
		"12658134172322624512087895071623532750001519942960405554845650333368647599808",
		"20763961938698988976208417731741058837854365735968696902157420039995071190051",
		"15629151274500299028863537382535241769408973898993482555097442222369536298769",
		"3302930679427685271940333650988715879176678940073974060817316157399962482951"}

	assignment.Siblings = [160]frontend.Variable{}
	for i := 0; i < 160; i++ {
		if i < len(encSiblings) {
			assignment.Siblings[i], _ = new(big.Int).SetString(encSiblings[i], 10)
		} else {
			assignment.Siblings[i] = big.NewInt(0)
		}
	}

	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))
}

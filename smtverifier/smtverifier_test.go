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

	assignment.Root, _ = new(big.Int).SetString("10231630067732032418299474999300673818856199563293770402682273793892295615796", 10)
	assignment.Key, _ = new(big.Int).SetString("873432238408170128747103711248787244651366455432", 10)
	assignment.Value, _ = new(big.Int).SetString("20", 10)

	encSiblings := []string{
		"16383576447159215769768642276201845720106009443497080943319659824221504437368",
		"9397881409034407145131092276609173661986975282781169183967420942140608554044",
		"21710180721102647226012023109254827809946922787497702964149388737548790104955",
		"6653039591693573146281004445905575758563705014333796295763984999642536123739",
		"19637833625829568397925201268938784934573906046500058722319740143210771947306",
		"9587778375427869868041036859495075545701145542689588716769811151013270824373",
		"8433100020982617748143688091368833564790102869936133347933857824100901970199"}

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

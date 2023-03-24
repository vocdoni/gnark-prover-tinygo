package circuit

import (
	"bytes"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"

	qt "github.com/frankban/quicktest"
)

type testingCircuit struct {
	A, B   frontend.Variable `gnark:",public"`
	Result frontend.Variable
}

func (c *testingCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(api.Mul(c.A, c.B), c.Result)
	return nil
}

func TestEnd2EndCircuit(t *testing.T) {
	c := qt.New(t)

	var circuit testingCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &circuit)
	c.Assert(err, qt.IsNil)

	var ccsBuff bytes.Buffer
	_, err = ccs.WriteTo(&ccsBuff)
	c.Assert(err, qt.IsNil)

	srs, err := test.NewKZGSRS(ccs)
	c.Assert(err, qt.IsNil)

	srsBuff := bytes.Buffer{}
	_, err = srs.WriteTo(&srsBuff)
	c.Assert(err, qt.IsNil)

	inputs := &testingCircuit{
		A:      2,
		B:      3,
		Result: 6,
	}
	wtns, err := frontend.NewWitness(inputs, ecc.BN254.ScalarField())
	c.Assert(err, qt.IsNil)

	wtnsBuff := bytes.Buffer{}
	_, err = wtns.WriteTo(&wtnsBuff)
	c.Assert(err, qt.IsNil)

	vKey, proof, pubWitness, err := GenerateProof(ccsBuff.Bytes(), srsBuff.Bytes(), wtnsBuff.Bytes())
	c.Assert(err, qt.IsNil)

	err = VerifyProof(srsBuff.Bytes(), vKey, proof, pubWitness)
	c.Assert(err, qt.IsNil)
}

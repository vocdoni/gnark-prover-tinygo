package groth16

import (
	"bytes"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

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
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	c.Assert(err, qt.IsNil)

	var ccsBuff bytes.Buffer
	_, err = ccs.WriteTo(&ccsBuff)
	c.Assert(err, qt.IsNil)

	provingKey, _, err := groth16.Setup(ccs)
	c.Assert(err, qt.IsNil)
	var pkeyBuff bytes.Buffer
	_, err = provingKey.WriteTo(&pkeyBuff)
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

	_, _, err = GenerateProof(ccsBuff.Bytes(), pkeyBuff.Bytes(), wtnsBuff.Bytes())
	c.Assert(err, qt.IsNil)
}

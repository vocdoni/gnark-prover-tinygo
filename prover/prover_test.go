package prover

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/kzg"
	"github.com/consensys/gnark/constraint"

	plonk "github.com/consensys/gnark/backend/plonk/bn254"
	cs "github.com/consensys/gnark/constraint/bn254"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"

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

func newKZGSRS(ccs constraint.ConstraintSystem) (*kzg.SRS, error) {
	nbConstraints := ccs.GetNbConstraints()
	sizeSystem := nbConstraints + ccs.GetNbPublicVariables()
	kzgSize := ecc.NextPowerOfTwo(uint64(sizeSystem)) + 3

	curveID := ecc.BN254
	alpha, err := rand.Int(rand.Reader, curveID.ScalarField())
	if err != nil {
		return nil, err
	}
	return kzg.NewSRS(kzgSize, alpha)
}

func TestEnd2EndCircuit(t *testing.T) {
	c := qt.New(t)

	var circuit testingCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &circuit)
	c.Assert(err, qt.IsNil)

	var ccsBuff bytes.Buffer
	_, err = ccs.WriteTo(&ccsBuff)
	c.Assert(err, qt.IsNil)

	srs, err := newKZGSRS(ccs)
	c.Assert(err, qt.IsNil)
	var srsBuff bytes.Buffer
	_, err = srs.WriteTo(&srsBuff)
	c.Assert(err, qt.IsNil)

	_scs := ccs.(*cs.SparseR1CS)
	provingKey, _, err := plonk.Setup(_scs, srs)
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

	_, _, err = GenerateProofPlonk(ccsBuff.Bytes(), srsBuff.Bytes(), pkeyBuff.Bytes(), wtnsBuff.Bytes())
	c.Assert(err, qt.IsNil)
}

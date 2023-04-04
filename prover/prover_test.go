package prover

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/vocdoni/gnark-crypto-bn254/ecc"
	"github.com/vocdoni/gnark-crypto-bn254/ecc/bn254/fr/kzg"
	"github.com/vocdoni/gnark-wasm-prover/constraint"
	cs "github.com/vocdoni/gnark-wasm-prover/csbn254"
	"github.com/vocdoni/gnark-wasm-prover/frontend"
	"github.com/vocdoni/gnark-wasm-prover/frontend/cs/scs"
	plonk "github.com/vocdoni/gnark-wasm-prover/prover"
	"github.com/vocdoni/gnark-wasm-prover/utils"

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

	curveID := utils.FieldToCurve(ccs.Field())
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

	_, _, err = GenerateProof(ccsBuff.Bytes(), srsBuff.Bytes(), pkeyBuff.Bytes(), wtnsBuff.Bytes())
	c.Assert(err, qt.IsNil)
}

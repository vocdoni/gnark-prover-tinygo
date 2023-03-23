package zkaddress

import (
	"gnark-prover-tinygo/internal/zkaddress"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

type testZkAddressCiruit struct {
	PrivateKey frontend.Variable
	PublicKey  frontend.Variable `gnark:",public"`
	Scalar     frontend.Variable `gnark:",public"`
}

func (circuit *testZkAddressCiruit) Define(api frontend.API) error {
	zkAddress, err := FromPrivate(api, circuit.PrivateKey)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.PrivateKey, zkAddress.Private)
	api.AssertIsEqual(circuit.PublicKey, zkAddress.Public)
	api.AssertIsEqual(circuit.Scalar, zkAddress.Scalar)
	return nil
}

func TestZkAddress(t *testing.T) {
	assert := test.NewAssert(t)

	seed := []byte("1b505cdafb4b1150b1a740633af41e5e1f19a5c4")
	zkAddr, err := zkaddress.FromBytes(seed)
	if err != nil {
		t.Error(err)
	}

	assignment := testZkAddressCiruit{
		PublicKey:  zkAddr.Public,
		PrivateKey: zkAddr.Private,
		Scalar:     zkAddr.Scalar,
	}

	var circuit testZkAddressCiruit
	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
}

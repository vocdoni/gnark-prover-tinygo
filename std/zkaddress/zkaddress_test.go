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
	Private frontend.Variable
	Public  frontend.Variable `gnark:",public"`
	Scalar  frontend.Variable `gnark:",public"`
}

func (circuit *testZkAddressCiruit) Define(api frontend.API) error {
	zkAddress, err := FromPrivate(api, circuit.Private)
	if err != nil {
		return err
	}

	api.Println("zkAddress", zkAddress.Private)
	api.Println("zkAddress", zkAddress.Public)
	api.Println("zkAddress", zkAddress.Scalar)
	api.Println("circuit", circuit.Private)
	api.Println("circuit", circuit.Public)
	api.Println("circuit", circuit.Scalar)

	api.AssertIsEqual(circuit.Private, zkAddress.Private)
	api.AssertIsEqual(circuit.Public, zkAddress.Public)
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
		Public:  zkAddr.Public,
		Private: zkAddr.Private,
		Scalar:  zkAddr.Scalar,
	}

	var circuit testZkAddressCiruit
	assert.SolvingSucceeded(&circuit, &assignment, test.WithCurves(ecc.BN254), test.WithBackends(backend.PLONK))
}

package zkaddress

import (
	"gnark-prover-tinygo/std/hash/poseidon"

	ecc "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
)

const DefaultZkAddressLen = 20

type ZkAddress struct {
	Private frontend.Variable
	Public  frontend.Variable
	Scalar  frontend.Variable
}

func FromPrivate(api frontend.API, private frontend.Variable) (ZkAddress, error) {
	curve, err := twistededwards.NewEdCurve(api, ecc.BN254)
	if err != nil {
		return ZkAddress{}, err
	}

	base := twistededwards.Point{
		X: curve.Params().Base[0],
		Y: curve.Params().Base[1],
	}
	point := curve.ScalarMul(base, private)
	public := poseidon.Hash(api, point.X, point.Y)

	bPublic := api.ToBinary(public, api.Compiler().FieldBitLen())
	scalar := api.FromBinary(bPublic[:DefaultZkAddressLen*8]...)

	return ZkAddress{
		Private: private,
		Public:  public,
		Scalar:  scalar,
	}, nil
}

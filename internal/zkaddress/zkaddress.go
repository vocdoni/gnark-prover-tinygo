package zkaddress

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

const DefaultZkAddressLen = 20

type ZkAddress struct {
	Private *big.Int
	Public  *big.Int
	Scalar  *big.Int
}

// func(zkAddr *ZkAddress) ArboKey(weight *big.Int) ([]byte, error) {
// 	return poseidon.Hash()
// }

func FromBytes(seed []byte) (*ZkAddress, error) {
	// Setup the curve
	c := twistededwards.GetEdwardsCurve()
	// Get scalar private key hashing the seed with Poseidon hash
	private, err := poseidon.HashBytes(seed)
	if err != nil {
		return nil, err
	}
	// Get the point of the curve that represents the public key multipliying
	// the private key scalar by the base of the curve
	point := new(twistededwards.PointAffine).ScalarMultiplication(&c.Base, private)
	// Get the single scalar that represents the publick key hashing X, Y point
	// coordenates with Poseidon hash
	bX, bY := new(big.Int), new(big.Int)
	bX = point.X.BigInt(bX)
	bY = point.Y.BigInt(bY)
	public, err := poseidon.Hash([]*big.Int{bX, bY})
	if err != nil {
		return nil, err
	}

	// truncate the most significant n bytes of the public key (little endian)
	// where n is the default ZkAddress length
	publicBytes := public.Bytes()
	m := len(publicBytes) - DefaultZkAddressLen
	scalar := new(big.Int).SetBytes(publicBytes[m:])
	return &ZkAddress{private, public, scalar}, nil
}

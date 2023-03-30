package zkaddress

import (
	"gnark-prover-tinygo/internal/arbo"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
)

const DefaultZkAddressLen = 20

type ZkAddress struct {
	Private *big.Int
	Public  *big.Int
	Scalar  *big.Int
}

func (zkAddr *ZkAddress) ArboBytes() []byte {
	return arbo.BigIntToBytes(DefaultZkAddressLen, zkAddr.Scalar)
}

func FromBytes(seed []byte) (*ZkAddress, error) {
	// Setup the curve
	c := twistededwards.GetEdwardsCurve()
	// Get scalar private key hashing the seed with Poseidon hash
	var hash arbo.HashMiMC
	bPrivate, err := hash.Hash(seed)
	if err != nil {
		return nil, err
	}
	private := arbo.BytesToBigInt(bPrivate)
	// Get the point of the curve that represents the public key multipliying
	// the private key scalar by the base of the curve
	point := new(twistededwards.PointAffine).ScalarMultiplication(&c.Base, private)
	// Get the single scalar that represents the publick key hashing X, Y point
	// coordenates with Poseidon hash
	bX := arbo.BigIntToBytes(arbo.HashFunctionMiMC.Len(), point.X.BigInt(new(big.Int)))
	bY := arbo.BigIntToBytes(arbo.HashFunctionMiMC.Len(), point.Y.BigInt(new(big.Int)))
	publicBytes, err := hash.Hash(bX, bY)
	if err != nil {
		return nil, err
	}

	// truncate the most significant n bytes of the public key (little endian)
	// where n is the default ZkAddress length
	scalar := arbo.BytesToBigInt(publicBytes[:DefaultZkAddressLen])
	return &ZkAddress{private, arbo.BytesToBigInt(publicBytes), scalar}, nil
}

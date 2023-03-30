package arbo

import (
	"encoding/hex"
	"math/big"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHashSha256(t *testing.T) {
	// Sha256 hash
	hashFunc := &HashSha256{}
	b := []byte("test")
	h, err := hashFunc.Hash(b)
	if err != nil {
		t.Fatal(err)
	}
	c := qt.New(t)
	c.Assert(hex.EncodeToString(h),
		qt.Equals,
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
}

func TestHashPoseidon(t *testing.T) {
	// Poseidon hash
	hashFunc := &HashPoseidon{}
	bLen := hashFunc.Len()
	h, err := hashFunc.Hash(
		BigIntToBytes(bLen, big.NewInt(1)),
		BigIntToBytes(bLen, big.NewInt(2)))
	if err != nil {
		t.Fatal(err)
	}
	hBI := BytesToBigInt(h)
	// value checked with circomlib
	c := qt.New(t)
	c.Assert(hBI.String(),
		qt.Equals,
		"7853200120776062878684798364095072458815029376092732009249414926327459813530")
}

func TestHashBlake2b(t *testing.T) {
	// Blake2b hash
	hashFunc := &HashBlake2b{}
	b := []byte("test")
	h, err := hashFunc.Hash(b)
	if err != nil {
		t.Fatal(err)
	}
	c := qt.New(t)
	c.Assert(hex.EncodeToString(h),
		qt.Equals,
		"928b20366943e2afd11ebc0eae2e53a93bf177a4fcf35bcc64d503704e65e202")
}

func TestHashMiMC(t *testing.T) {
	c := qt.New(t)
	// MiMC hash
	hashFunc := &HashMiMC{}
	input, _ := new(big.Int).SetString("18045289051299654077710208499747278752099041449041972372412271818361923969579", 10)

	h, err := hashFunc.Hash(BigIntToBytes(hashFunc.Len(), input))
	c.Assert(err, qt.IsNil)
	c.Assert(BytesToBigInt(h).String(), qt.Equals, "3063534589717532550748064309090249843004878651469661670941107341968586307081")
}

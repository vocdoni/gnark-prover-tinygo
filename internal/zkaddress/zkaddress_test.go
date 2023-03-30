package zkaddress

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFromBytes(t *testing.T) {
	c := qt.New(t)

	seed := []byte("1b505cdafb4b1150b1a740633af41e5e1f19a5c4")
	zkAddr, err := FromBytes(seed)
	c.Assert(err, qt.IsNil)
	c.Assert(zkAddr.Private.String(), qt.Equals, "20104241803663641422577121134203490505137011783614913652735802145961801733870")
	c.Assert(zkAddr.Public.String(), qt.Equals, "2493779843424947948760282772832914324283078143588187307135787195808806220423")
	c.Assert(zkAddr.Scalar.String(), qt.Equals, "778541079330801545513944229279598209414021919367")
}

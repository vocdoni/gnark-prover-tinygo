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
	c.Assert(zkAddr.Private.String(), qt.Equals, "4942627222315175708338446333266917590203114611189826202299682977805004266061")
	c.Assert(zkAddr.Public.String(), qt.Equals, "3560047183729315007179763679233971090993938368129267583816767417160809903594")
	c.Assert(zkAddr.Scalar.String(), qt.Equals, "70956925796045393021907218364509574825414292970")
}

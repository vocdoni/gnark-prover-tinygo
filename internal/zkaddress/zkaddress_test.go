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
	c.Assert(zkAddr.Private.String(), qt.Equals, "12007696602022466067210558438468234995085206818257350359618361229442198701667")
	c.Assert(zkAddr.Public.String(), qt.Equals, "19597797733822453932297698102694210740977986979020091017779598307769964166976")
	c.Assert(zkAddr.Scalar.String(), qt.Equals, "647903732945896451226807429503635300036365909824")
}

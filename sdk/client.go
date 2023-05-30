package sdk

import (
	"fmt"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

func ComputexGxH(x kyber.Scalar, G, H kyber.Point) (kyber.Point, kyber.Point) {
	suite := edwards25519.NewBlakeSHA256Ed25519()

	xG := suite.Point().Mul(x, G)
	xH := suite.Point().Mul(x, H)

	fmt.Printf("xG:\t%s\n xH\t%s\n\n", xG, xH)

	return xG, xH
}

func ComputevGvH(G, H kyber.Point, xG, xH kyber.Point) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	// unknown
	v := suite.Scalar().Pick(suite.RandomStream())
	vG := suite.Point().Mul(v, G)
	vH := suite.Point().Mul(v, H)

	// Bob can now generate a challenge, without requiring Alice to send it (non-interactive).

	//  For this Bob, takes a hash of xG, xH, vG and vH:
	h := suite.Hash()
	xG.MarshalTo(h)
	xH.MarshalTo(h)
	vG.MarshalTo(h)
	vH.MarshalTo(h)
	cb := h.Sum(nil)
	c := suite.Scalar().Pick(suite.XOF(cb))
}

func ComputeRrGrH() {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	
	// Bob can then compute r, rG, and rH with:
	// Response
	r := suite.Scalar()
	r.Mul(x, c).Sub(v, r)

	rG := suite.Point().Mul(r, G)
	rH := suite.Point().Mul(r, H)
	cxG := suite.Point().Mul(c, xG)
	cxH := suite.Point().Mul(c, xH)
	a := suite.Point().Add(rG, cxG)
	b := suite.Point().Add(rH, cxH)

	fmt.Printf("Alice sends challenge:\n c: %s\n\n", c)
	fmt.Printf("Bob computes:\n v:\t%s\n r:\t%s\n\n", v, r)

	if !(vG.Equal(a) && vH.Equal(b)) {
		fmt.Printf("Incorrect proof!\n")
	} else {
		fmt.Printf("Proof correct\n")
	}
}

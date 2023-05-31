package sdk

import (
	"crypto/sha256"
	"io"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/util/random"
)

var rng = random.New()

// 1 stage - server
func ComputeX(m []byte, suite *edwards25519.SuiteEd25519) kyber.Scalar {
	// In Go, we can convert our password (m) into a value by taking a hash of it (SHA-256):
	message := []byte(m)
	scal := sha256.Sum256(message[:])

	x := suite.Scalar().SetBytes(scal[:32])

	return x
}

// 2 stage - server
func ComputeGH(suite *edwards25519.SuiteEd25519) (G, H kyber.Point) {
	// Next we can generate our elliptic curve points (G and H):
	G = suite.Point().Pick(rng)
	H = suite.Point().Pick(rng)

	return G, H
}

func ComputexGxH(suite *edwards25519.SuiteEd25519, G, H kyber.Point, x kyber.Scalar) (xG kyber.Point, xH kyber.Point) {
	xG = suite.Point().Mul(x, G)
	xH = suite.Point().Mul(x, H)
	return xG, xH
}

func ComputeVvGvH(suite *edwards25519.SuiteEd25519, G, H kyber.Point) (kyber.Scalar, kyber.Point, kyber.Point) {
	v := suite.Scalar().Pick(suite.RandomStream())
	vG := suite.Point().Mul(v, G)
	vH := suite.Point().Mul(v, H)
	return v, vG, vH
}

func ComputeC(suite *edwards25519.SuiteEd25519, G, H, xG, xH kyber.Point) kyber.Scalar {
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

	return c
}

func ComputeRrGrH(suite *edwards25519.SuiteEd25519, c, x, v kyber.Scalar, G, H kyber.Point) (r kyber.Scalar, rG, rH kyber.Point) {
	r = suite.Scalar()
	r.Mul(x, c).Sub(v, r)

	rG = suite.Point().Mul(r, G)
	rH = suite.Point().Mul(r, H)

	return r, rG, rH
}

func ComputeAB(suite *edwards25519.SuiteEd25519, c kyber.Scalar, xH, xG, rG, rH kyber.Point) (kyber.Point, kyber.Point) {
	cxG := suite.Point().Mul(c, xG)
	cxH := suite.Point().Mul(c, xH)
	a := suite.Point().Add(rG, cxG)
	b := suite.Point().Add(rH, cxH)
	return a, b
}

func Valid(vG, vH, a, b kyber.Point) bool {
	if !(vG.Equal(a) && vH.Equal(b)) {
		return false
	} else {
		return true
	}
}

// GenerateKey is used to generate a new key
func GenerateKey(reader io.Reader, count int) ([]byte, error) {
	// Generate a 256bit key
	buf := make([]byte, count)
	_, err := reader.Read(buf)

	return buf, err
}

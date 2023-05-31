package sdk

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"go.dedis.ch/kyber/v3/group/edwards25519"
)

func TestFiatShamir(t *testing.T) {
	m, _ := GenerateKey(rand.Reader, 32)

	suite := edwards25519.NewBlakeSHA256Ed25519()

	x := ComputeX(m, suite)  // Alice and Bob
	G, H := ComputeGH(suite) // (Alice) Next we can generate our elliptic curve points (G and H)

	fmt.Printf("Bob and Alice agree:\n G:\t%s\n H:\t%s\n\n", G, H)
	fmt.Printf("Bob's Password:\t%s\n", hex.EncodeToString(m))
	fmt.Printf("Bob's Secret (x):\t%s\n\n", x)

	xG, xH := ComputexGxH(suite, G, H, x) // (Bob) The values passed to Alice can be generated with

	fmt.Printf("Bob sends these values:\n xG:\t%s\n xH\t%s\n\n", xG, xH)

	c := ComputeC(suite, G, H, xG, xH) // (Bob) Bob can now generate a challenge, without requiring Alice to send it (non-interactive).

	v, vG, vH := ComputeVvGvH(suite, G, H) // unknown

	// Bob can then compute r, rG, and rH with:
	// Response
	r, rG, rH := ComputeRrGrH(suite, c, x, v, G, H)

	a, b := ComputeAB(suite, c, xH, xG, rG, rH) // (Alice) // Alice then checks

	fmt.Printf("Alice sends challenge:\n c: %s\n\n", c)
	fmt.Printf("Bob computes:\n v:\t%s\n r:\t%s\n\n", v, r)

	if !Valid(vG, vH, a, b) { // Alice
		fmt.Printf("Incorrect proof!\n")
	} else {
		fmt.Printf("Proof correct\n")
	}
}

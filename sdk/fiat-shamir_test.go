package sdk

import (
	"crypto/rand"
	"fmt"
	"testing"

	"go.dedis.ch/kyber/v3/group/edwards25519"
)

func TestFiatShamir(t *testing.T) {
	m, _ := GenerateKey(rand.Reader, 32) // (server)

	suite := edwards25519.NewBlakeSHA256Ed25519() // (client & server)
	x := ComputeX(m, suite)                       // (client)
	G, H := ComputeGH(suite)                      // (server) Next we can generate our elliptic curve points (G and H)

	xG, xH := ComputexGxH(suite, G, H, x)           // (client) The values passed to Alice can be generated with
	c := ComputeC(suite, G, H, xG, xH)              // (client) Bob can now generate a challenge, without requiring Alice to send it (non-interactive).
	v, vG, vH := ComputeVvGvH(suite, G, H)          // (client) unknown
	_, rG, rH := ComputeRrGrH(suite, c, x, v, G, H) // (client) Bob can then compute r, rG, and rH with:

	a, b := ComputeAB(suite, c, xH, xG, rG, rH) // (server) Alice then checks
	valid := Valid(vG, vH, a, b)                // (server)

	// check
	if valid {
		fmt.Printf("Incorrect proof!\n")
	} else {
		fmt.Printf("Proof correct\n")
	}
}

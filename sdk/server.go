package sdk

import (
	"crypto/sha256"
	"io"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/util/random"

	hclog "github.com/hashicorp/go-hclog"
)

var rng = random.New()

func ECCGH(m []byte, logger hclog.Logger) (kyber.Point, kyber.Point, error) {
	suite := edwards25519.NewBlakeSHA256Ed25519()

	// In Go, we can convert our password (m) into a value by taking a hash of it (SHA-256):
	message := []byte(m)
	scal := sha256.Sum256(message[:])

	x := suite.Scalar().SetBytes(scal[:32])

	// Next we can generate our elliptic curve points (G and H):
	G := suite.Point().Pick(rng)
	H := suite.Point().Pick(rng)

	logger.Info("G:\t%s\n H:\t%s\n\n", G, H)

	logger.Info("password:\t%s\n", m)
	logger.Info("secret (x):\t%s\n\n", x)

	return G, H, nil
}

// GenerateKey is used to generate a new key
func GenerateKey(reader io.Reader) ([]byte, error) {
	// Generate a 256bit key
	buf := make([]byte, 32)
	_, err := reader.Read(buf)

	return buf, err
}

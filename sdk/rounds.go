package sdk

import (
	"github.com/fxamacker/cbor/v2"
	"go.dedis.ch/kyber/v3"
)

type Round1 struct {
	Message []byte
	G       kyber.Point
	H       kyber.Point
}

type round1 struct {
	Message []byte
	G       kyber.Point
	H       kyber.Point
}

func (r1 *Round1) toMarshallable() *round1 {
	return &round1{
		Message: r1.Message,
		G:       r1.G,
		H:       r1.H,
	}
}

func (r1 *Round1) MarshalBinary() ([]byte, error) {
	return cbor.Marshal(r1.toMarshallable())
}

func (r1 *Round1) UnmarshalBinary(data []byte) error {
	if err := cbor.Unmarshal(data, r1); err != nil {
		return err
	}

	return nil
}

type Round2 struct {
	C kyber.Scalar

	XH kyber.Point
	XG kyber.Point

	RG kyber.Point
	RH kyber.Point

	VG kyber.Point
	VH kyber.Point
}

type round2 struct {
	C kyber.Scalar

	XH kyber.Point
	XG kyber.Point

	RG kyber.Point
	RH kyber.Point

	VG kyber.Point
	VH kyber.Point
}

func (r2 *Round2) toMarshallable() *round2 {
	return &round2{
		C: r2.C,

		XH: r2.XH,
		XG: r2.XG,

		RG: r2.RG,
		RH: r2.RH, 

		VG: r2.VG, 
		VH: r2.VH,
	}
}

func (r2 *Round2) MarshalBinary() ([]byte, error) {
	return cbor.Marshal(r2.toMarshallable())
}

func (r2 *Round2) UnmarshalBinary(data []byte) error {
	if err := cbor.Unmarshal(data, r2); err != nil {
		return err
	}

	return nil
}

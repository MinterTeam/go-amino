package data

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
)

// Encoder is a global setting for all byte encoding
// This is the default.  Please override in the main()/init()
// of your program to change how byte slices are presented
var (
	Encoder       ByteEncoder = hexEncoder{}
	HexEncoder                = hexEncoder{}
	B64Encoder                = base64Encoder{base64.URLEncoding}
	RawB64Encoder             = base64Encoder{base64.RawURLEncoding}
)

// Bytes is a special byte slice that allows us to control the
// serialization format per app.
//
// Thus, basecoin could use hex, another app base64, and a third
// app base58...
type Bytes []byte

func (b Bytes) MarshalJSON() ([]byte, error) {
	return Encoder.Marshal(b)
}

func (b *Bytes) UnmarshalJSON(data []byte) error {
	ref := (*[]byte)(b)
	return Encoder.Unmarshal(ref, data)
}

// ByteEncoder handles both the marshalling and unmarshalling of
// an arbitrary byte slice.
//
// All Bytes use the global Encoder set in this package.
// If you want to use this encoding for byte arrays, you can just
// implement a simple custom marshaller for your byte array
//
//   type Dings [64]byte
//
//   func (d Dings) MarshalJSON() ([]byte, error) {
//     return data.Encoder.Marshal(d[:])
//   }
//
//   func (d *Dings) UnmarshalJSON(data []byte) error {
//     ref := (*d)[:]
//     return data.Encoder.Unmarshal(&ref, data)
//   }
type ByteEncoder interface {
	Marshal(bytes []byte) ([]byte, error)
	Unmarshal(dst *[]byte, src []byte) error
}

// hexEncoder implements ByteEncoder encoding the slice as a hexidecimal
// string
type hexEncoder struct{}

func (h hexEncoder) _assertByteEncoder() ByteEncoder {
	return h
}

func (_ hexEncoder) Unmarshal(dst *[]byte, src []byte) (err error) {
	var s string
	err = json.Unmarshal(src, &s)
	if err != nil {
		return errors.Wrap(err, "parse string")
	}
	// and interpret that string as hex
	*dst, err = hex.DecodeString(s)
	return err
}

func (_ hexEncoder) Marshal(bytes []byte) ([]byte, error) {
	s := hex.EncodeToString(bytes)
	return json.Marshal(s)
}

// base64Encoder implements ByteEncoder encoding the slice as
// base64 url-safe encoding
type base64Encoder struct {
	*base64.Encoding
}

func (e base64Encoder) _assertByteEncoder() ByteEncoder {
	return e
}

func (e base64Encoder) Unmarshal(dst *[]byte, src []byte) (err error) {
	var s string
	err = json.Unmarshal(src, &s)
	if err != nil {
		return errors.Wrap(err, "parse string")
	}
	*dst, err = e.DecodeString(s)
	return err
}

func (e base64Encoder) Marshal(bytes []byte) ([]byte, error) {
	s := e.EncodeToString(bytes)
	return json.Marshal(s)
}

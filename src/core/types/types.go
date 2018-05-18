package types

/*
import (
	"crypto/sha1"
	"encoding/base64"
)
*/

//import "encoding"

const (
	// HashLength - basic, arbitrary specification for later hash format changes
	HashLength = 32

	// AddressLength - basic, arbitrary specification for later address format changes
	AddressLength = 20
)

//Hash - arbitraty data hash type
type Hash [HashLength]byte

//Address - address of specificed byte length
type Address [AddressLength]byte

// BytesToAddress - Set address instance to byte array.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// SetBytes - Sets the address to the value of b. If b is larger than len(a) it will panic
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

package types

import (
	"encoding/hex"
	"math/big"
)

const (
	// HashLength - basic, arbitrary specification for later hash format changes
	HashLength = 32

	// AddressLength - basic, arbitrary specification for later address format changes
	AddressLength = 20
)

// Hash - arbitraty data hash type
type Hash [HashLength]byte

// Address - address of specificed byte length
type Address [AddressLength]byte

// Identifier - Byte array representing identification of dataset
type Identifier []byte

// Weight - data representing computational power for one transaction
type Weight *int

// URL - API available reference to network account
type URL struct {
	Scheme string
	Path   string
}

// BytesToAddress - Set address instance to byte array.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// SetBytes - Sets the address to the value of b.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// BytesToHash - Set hash instance to byte array.
func BytesToHash(b []byte) Hash {
	var a Hash
	copy(a[:], b)
	return a
}

// SetBytes - Sets the address to the value of b.
func (a *Hash) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b):]
	}
	copy(a[len(b):], b)
}

// BigToAddress - Convert int to Address
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress - Convert hex string to Address
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// FromHex - Generate byte array from hex string
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes - convert Hex string to byte array
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)

	return h
}

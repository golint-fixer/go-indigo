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
type Weight *big.Int

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

// BigToAddress - Convert big.Int to Address
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

// BigToWeight - Convert specified big.Int value to instance of weight struct
func BigToWeight(val *big.Int) Weight {
	var weightv Weight
	weightv = val
	return weightv
}

// IntToWeight - Convert specified Int value to instance of weight struct
func IntToWeight(val int) Weight {
	var weightv Weight
	weightv = big.NewInt(int64(val))
	return weightv
}

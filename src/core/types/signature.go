package types

import (
	"encoding/hex"
	"math/big"
)

//Signature - data representing digital verification, as well as any payload attatched to the verification.
type Signature []byte

// BytesToSignature - Convert specified byte array to tx signature.
func BytesToSignature(b []byte) Signature {
	var a Signature
	a = b
	return a
}

// BigToSignature - Convert big.Int to Signature
func BigToSignature(b *big.Int) Signature { return BytesToSignature(b.Bytes()) }

// HexToSignature - Convert hex string to Signature
func HexToSignature(s string) Signature { return BytesToSignature(fromHex(s)) }

// fromHex - Generate byte array from hex string
func fromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return hex2Bytes(s)
}

// hex2Bytes - convert Hex string to byte array
func hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)

	return h
}

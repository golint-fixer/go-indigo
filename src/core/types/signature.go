package types

import "math/big"

//Signature - data representing digital verification, as well as any payload attatched to the verification.
type Signature []byte

// BytesToSignature - Convert specified byte array to tx signature.
func BytesToSignature(b []byte) Signature {
	var a Signature
	a = b
	return a
}

// BigToSignature - Convert int to Signature
func BigToSignature(b *big.Int) Signature { return BytesToSignature(b.Bytes()) }

// IntToSignature - Convert int to Signature
func IntToSignature(b *int) Signature { return BytesToSignature((*big.NewInt(int64(*b))).Bytes()) }

// HexToSignature - Convert hex string to Signature
func HexToSignature(s string) Signature { return BytesToSignature(FromHex(s)) }

// Block components are based on go-ethereum/common/types.go

package common

import (
	"encoding/hex"
)

const (
	HashLength    = 32
	AddressLength = 20
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

type Address [AddressLength]byte

func (h Hash) ToHex() string {
	var b = make([]byte, HashLength)
	for i:=0;i<HashLength;i++{
		b[i] = h[i]
	}

	hex := Bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}

func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

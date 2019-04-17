// Block components are based on go-ethereum/common/types.go

package common

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/altair-lab/xoreum/common/math"
)

const (
	HashLength    = 32
	AddressLength = 32 // can be changed later
)

var (
	Difficulty = math.BigPow(2, 256-1) // mining difficulty: 10
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

type Address [AddressLength]byte

func (h Hash) ToBigInt() *big.Int {
	byteArr := []byte{}

	for i := 0; i < HashLength; i++ {
		byteArr = append(byteArr, h[i])
	}

	r := new(big.Int)
	r.SetBytes(byteArr)
	return r
}

func (h Hash) ToHex() string {
	var b = make([]byte, HashLength)
	for i := 0; i < HashLength; i++ {
		b[i] = h[i]
	}

	hex := Bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}

func (a Address) ToHex() string {
	var b = make([]byte, AddressLength)
	for i := 0; i < AddressLength; i++ {
		b[i] = a[i]
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

func HexToHash(s string) Hash {
	return Hash{}
}

func ToBytes(v interface{}) []byte {
	return []byte(fmt.Sprintf("%v", v))
}

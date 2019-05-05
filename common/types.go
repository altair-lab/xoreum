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
	// original
	Difficulty = math.BigPow(2, 256-1) // mining difficulty: 10

	// this is for test
	//Difficulty = math.BigPow(2, 260)
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

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

// CopyBytes returns an exact copy of the provided bytes.
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// Num2Bytes returns a byte array representation of a uint64 value
func Num2Bytes(num uint64) []byte {
	h := fmt.Sprintf("%x", num)
	bytes, _ := hex.DecodeString(h)
	return bytes
}

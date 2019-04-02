// Block components are based on go-ethereum/common/types.go

package common

import "fmt"

const (
	HashLength		= 32
	AddressLength	= 20
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

type Address [AddressLength]byte

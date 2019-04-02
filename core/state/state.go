package state

import (
	"github.com/altair-lab/xoreum/common"
)

type State map[common.Address]uint64

type Account struct {
	Address common.Address
	Nonce   uint64
	Balance uint64
}

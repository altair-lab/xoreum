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

func NewAccount(address common.Address, nonce uint64, balance uint64) *Account {
	return newAccount(address, nonce, balance)
}

func newAccount(address common.Address, nonce uint64, balance uint64) *Account {
	return &Account{
		Address: address,
		Nonce:   nonce,
		Balance: balance,
	}
}

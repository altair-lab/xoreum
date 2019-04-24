package state

import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
)


type State map[common.Address]*Account

type Account struct {
	Address common.Address
	Nonce   uint64
	Balance uint64
}

func NewState() State {
	return State{}
}

func (s State) Add(acc *Account) {
	s[acc.Address] = acc
}

func (s State) Print() {
	for _, v := range s {
		v.Print()
	}
}

func (s State) GetBalance(address common.Address) uint64 {
	return s[address].Balance
}

func (s State) GetNonce(address common.Address) uint64 {
	return s[address].Nonce
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

func (acc *Account) Print() {
	fmt.Printf("Address: %x   Nonce: %d   Balance: %d\n", acc.Address, acc.Nonce, acc.Balance)
}

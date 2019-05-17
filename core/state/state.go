package state

import (
	"crypto/ecdsa"
	"fmt"
)

type State map[ecdsa.PublicKey]*Account

type Account struct {
	//Address common.Address
	PublicKey *ecdsa.PublicKey
	Nonce     uint64
	Balance   uint64
}

func NewState() State {
	return State{}
}

func (s State) Add(acc *Account) {
	s[*acc.PublicKey] = acc
}

func (s State) Print() {
	sum := uint64(0)
	for _, v := range s {
		v.PrintAccount()
		sum += v.Balance
	}
	fmt.Println("balance sum:", sum)
}

func (s State) GetBalance(pubkey *ecdsa.PublicKey) uint64 {
	return s[*pubkey].Balance
}

func (s State) GetNonce(pubkey *ecdsa.PublicKey) uint64 {
	return s[*pubkey].Nonce
}

func (s State) NewAccount(pubkey *ecdsa.PublicKey, nonce uint64, balance uint64) *Account {
	acc := newAccount(pubkey, nonce, balance)
	s.Add(acc)
	return acc
}

func NewAccount(pubkey *ecdsa.PublicKey, nonce uint64, balance uint64) *Account {
	return newAccount(pubkey, nonce, balance)
}

func newAccount(pubkey *ecdsa.PublicKey, nonce uint64, balance uint64) *Account {
	return &Account{
		PublicKey: pubkey,
		Nonce:     nonce,
		Balance:   balance,
	}
}

func (acc *Account) Print() {
	fmt.Printf("PublicKey: %x   Nonce: %d   Balance: %d\n", acc.PublicKey, acc.Nonce, acc.Balance)
}

func (acc *Account) PrintAccount() {
	fmt.Println("publickey:", acc.PublicKey, "/ nonce:", acc.Nonce, "/ balance:", acc.Balance)
}

func (acc *Account) Copy() *Account {
	return NewAccount(acc.PublicKey, acc.Nonce, acc.Balance)
}

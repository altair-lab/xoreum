package state

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/altair-lab/xoreum/common"
	//"github.com/altair-lab/xoreum/common"
)

type State map[ecdsa.PublicKey]common.Hash // pubkey - TxHash (user's current tx hash)

type Accounts map[ecdsa.PublicKey]*Account // pubkey - Account

type Account struct {
	PublicKey *ecdsa.PublicKey
	Nonce     uint64
	Balance   uint64
}

func NewAccounts() Accounts {
	return Accounts{}
}

func (s Accounts) Add(acc *Account) {
	s[*acc.PublicKey] = acc
}

func (s Accounts) Print() {
	sum := uint64(0)
	for _, v := range s {
		v.PrintAccount()
		sum += v.Balance
	}
	fmt.Println("balance sum:", sum)
}

func (s Accounts) PrintAccountsSum() {
	sum := uint64(0)
	for _, v := range s {
		sum += v.Balance
	}
	fmt.Println("accounts balance sum:", sum)
}

func (s Accounts) GetBalance(pubkey *ecdsa.PublicKey) uint64 {
	return s[*pubkey].Balance
}

func (s Accounts) GetNonce(pubkey *ecdsa.PublicKey) uint64 {
	return s[*pubkey].Nonce
}

func (s Accounts) NewAccount(pubkey *ecdsa.PublicKey, nonce uint64, balance uint64) *Account {
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
	fmt.Println("publickey:", acc.PublicKey /*"publickey.Curve.params():", acc.PublicKey.Curve.Params(),*/, "/ nonce:", acc.Nonce, "/ balance:", acc.Balance)
}

func (acc *Account) Copy() *Account {
	return NewAccount(acc.PublicKey, acc.Nonce, acc.Balance)
}

func (s State) Print() {
	for k, v := range s {
		fmt.Println("pubkey:", k, "\n\t/ txhash:", v.ToHex())
	}
}

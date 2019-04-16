package types

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/crypto"
)

type Transaction struct {
	data txdata
	hash common.Hash

	// signature values
	Sender_R    *big.Int
	Sender_S    *big.Int
	Recipient_R *big.Int
	Recipient_S *big.Int
}

// Transactions is a Transaction slice type for basic sorting
type Transactions []*Transaction

/*
type txdata struct {
	Participants      []*common.Address
	ParticipantNonces []uint64
	XORs              []uint64
	Payload           []byte

	// Signature values
}
*/

// simple implementation
type txdata struct {
	AccountNonce uint64
	Sender    *ecdsa.PublicKey
	Recipient *ecdsa.PublicKey
	Amount    uint64
}

func NewTransaction(from ecdsa.PublicKey, to ecdsa.PublicKey, amount uint64) *Transaction {
	return newTransaction(&from, &to, amount)
}


// [TODO] Set nonce
func newTransaction(from *ecdsa.PublicKey, to *ecdsa.PublicKey, amount uint64) *Transaction {
	nonce := uint64(0)

	d := txdata{
		AccountNonce: nonce,
		Sender:       from,
		Recipient:    to,
		Amount:       amount,
	}

	return &Transaction{data: d}
}

func (tx *Transaction) Nonce() uint64	{ return tx.data.AccountNonce }
func (tx *Transaction) Value() uint64 { return tx.data.Amount } 
func (tx *Transaction) Sender() common.Address { return *tx.data.Sender } // Temporal function until signature is implemented
func (tx *Transaction) Recipient() common.Address { return *tx.data.Recipient }

func (tx *Transaction) Hash() common.Hash {
	return crypto.Keccak256Hash(common.ToBytes(*tx))
}

func (txs *Transactions) Hash() common.Hash {
  return crypto.Keccak256Hash(common.ToBytes(*txs))
}

func (tx *Transaction) GetTxdataHash() []byte {
	return crypto.Keccak256(common.ToBytes(tx.data))
}

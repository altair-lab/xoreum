package types

import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/crypto"
)

type Transaction struct {
	data txdata
	hash common.Hash
}

// Transactions is a Transaction slice type for basic sorting
type Transactions []*Transaction

// txdata could be generated between more than 2 participants
// For example, if A, B, C are participants, data of txdata is
// participants: [A, B, C]
// participantNonces: [10, 3, 5]
// XORs : ['1234', '3245', '4313']
// Payload : ""
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
	AccountNonce	uint64
	Sender		*common.Address // [TODO] Implement it using signature values (ref: transaction_signing.go)
	Recipient	*common.Address
	Amount		uint64
}


func NewTransaction(from common.Address, to common.Address, amount uint64) *Transaction {
	return newTransaction(&from, &to, amount)
}

func newTransaction(from *common.Address, to *common.Address, amount uint64) *Transaction {
	// [TODO] Set nonce
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
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("%v", *tx)))
}

func (txs *Transactions) Hash() common.Hash {
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("%v", *txs)))
}

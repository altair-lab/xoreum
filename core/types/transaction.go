package types

import (
	"github.com/altair-lab/xoreum/common"
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
	Sender		*common.Address
	Recipient	*common.Address
	Amount		uint64
}


func NewTransaction(from common.Address, to common.Address, amount uint64) *Transaction {
	return newTransaction(&from, &to, amount)
}

func newTransaction(from *common.Address, to *common.Address, amount uint64) *Transaction {
	nonce := uint64(0)

	d := txdata{
		AccountNonce: nonce,
		Sender:       from,
		Recipient:    to,
		Amount:       amount,
	}

	return &Transaction{data: d}
}

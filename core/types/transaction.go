package types

import (
	"crypto/ecdsa"
	"math/big"
	"github.com/altair-lab/xoreum/common"
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
	Sender       *ecdsa.PublicKey
	Recipient    *ecdsa.PublicKey
	Amount       uint64
}

func NewTransaction(nonce uint64, from ecdsa.PublicKey, to ecdsa.PublicKey, amount uint64) *Transaction {
	return newTransaction(nonce, &from, &to, amount)
}

func newTransaction(nonce uint64, from *ecdsa.PublicKey, to *ecdsa.PublicKey, amount uint64) *Transaction {
	d := txdata{
		AccountNonce: nonce,
		Sender:       from,
		Recipient:    to,
		Amount:       amount,
	}

	return &Transaction{data: d}
}

func (tx *Transaction) Nonce() uint64               { return tx.data.AccountNonce }
func (tx *Transaction) Value() uint64               { return tx.data.Amount }
func (tx *Transaction) Sender() *ecdsa.PublicKey    { return tx.data.Sender } // Temporal function until signature is implemented
func (tx *Transaction) Recipient() *ecdsa.PublicKey { return tx.data.Recipient }

func (tx *Transaction) Hash() common.Hash {
	return crypto.Keccak256Hash(common.ToBytes(*tx))
}

func (txs *Transactions) Hash() common.Hash {
	return crypto.Keccak256Hash(common.ToBytes(*txs))
}

func (tx *Transaction) GetTxdataHash() []byte {
	return crypto.Keccak256(common.ToBytes(tx.data))
}

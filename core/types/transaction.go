package types

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/crypto"
)

type Transaction struct {
	data txdata
	hash common.Hash

	// signature values (old version)
	/*Sender_R    *big.Int
	Sender_S    *big.Int
	Recipient_R *big.Int
	Recipient_S *big.Int*/

	// signature values of participants (new version)
	Signature_R []*big.Int
	Signature_S []*big.Int
}

// Transactions is a Transaction slice type for basic sorting
type Transactions []*Transaction

// simple implementation
type txdata struct {
	// old version fields
	AccountNonce uint64
	Sender       *ecdsa.PublicKey
	Recipient    *ecdsa.PublicKey
	Amount       uint64

	// new version fields
	Participants []*ecdsa.PublicKey
	PostStates   []*state.Account
	PrevTxHashes []*common.Hash
}

func NewTransaction(participants []*ecdsa.PublicKey, postStates []*state.Account, prevTxHashes []*common.Hash) *Transaction {
	d := txdata{
		Participants: participants,
		PostStates:   postStates,
		PrevTxHashes: prevTxHashes,
	}

	tx := Transaction{data: d}

	// dynamic allocation
	length := len(participants)
	tx.Signature_R = make([]*big.Int, length)
	tx.Signature_S = make([]*big.Int, length)

	return &tx
}

func (tx *Transaction) Nonce() uint64 { return tx.data.AccountNonce }

func (tx *Transaction) Value() uint64 { return tx.data.Amount }

func (tx *Transaction) Sender() ecdsa.PublicKey { return *tx.data.Sender } // Temporal function until signature is implemented

func (tx *Transaction) Recipient() ecdsa.PublicKey { return *tx.data.Recipient }

func (tx *Transaction) Hash() common.Hash {
	return crypto.Keccak256Hash(common.ToBytes(*tx))
}

// TODO: i think this should be changed
func (txs *Transactions) Hash() common.Hash {
	return crypto.Keccak256Hash(common.ToBytes(*txs))
}

// insert tx into txs
func (txs *Transactions) Insert(tx *Transaction) {
	*txs = append(*txs, tx)
}

func (tx *Transaction) GetTxdataHash() []byte {
	return crypto.Keccak256(common.ToBytes(tx.data))
}

func (tx *Transaction) PrintTx() {
	for i := 0; i < len(tx.data.Participants); i++ {
		fmt.Println("participant ", i)
		fmt.Println("public key: ", tx.data.Participants[i])
		fmt.Println("post state: ", tx.data.PostStates[i])
		fmt.Println("previous tx hash: ", tx.data.PrevTxHashes[i])
	}
}

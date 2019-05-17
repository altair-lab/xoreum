package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/rlp"
)

var (
	// incorrect block's number (current_block_number + 1 != insert_block's_number)
	ErrDiffFieldLength = errors.New("tx has different field array length")

	ErrInvalidPostStates = errors.New("tx's PostStates's Account is not same with tx's Participants")

	ErrInvalidPrevTxHashes = errors.New("Account in tx's PrevTxHashes is not match with Participants")
)

type Transaction struct {
	Data Txdata      `json:"d"`
	Hash common.Hash `json:"h"`

	// signature values of participants
	Signature_R []*big.Int `json:"r"`
	Signature_S []*big.Int `json:"s"`
}

// Transactions is a Transaction slice type for basic sorting
type Transactions []*Transaction

// simple implementation
type Txdata struct {
	// new version fields
	Participants []*ecdsa.PublicKey `json:"participants"`
	PostStates   []*state.Account   `json:"poststates"`
	PrevTxHashes []*common.Hash     `json:"prevtxhashes"`
}

func NewTransaction(participants []*ecdsa.PublicKey, postStates []*state.Account, prevTxHashes []*common.Hash) *Transaction {
	d := Txdata{
		Participants: participants,
		PostStates:   postStates,
		PrevTxHashes: prevTxHashes,
	}

	tx := Transaction{Data: d}

	// dynamic allocation
	length := len(participants)
	tx.Signature_R = make([]*big.Int, length)
	tx.Signature_S = make([]*big.Int, length)
	tx.Hash = tx.GetHash()

	return &tx
}

func UnmarshalJSON(txbuf []byte) *Transaction {
	// Get Participants Length
	d := Transaction{}
	json.Unmarshal(txbuf, &d)
	length := len(d.Data.Participants)

	// [BUGFIX] tx.txdata.Participants[i].Curve == <nil>
	participants := make([]*(ecdsa.PublicKey), length)
	postStates := make([]*(state.Account), length)

	for i := 0; i < length; i++ {
		participants[i] = &ecdsa.PublicKey{Curve: &elliptic.CurveParams{}}
		postStates[i] = &state.Account{PublicKey: participants[i]}
	}
	txdata := Txdata{Participants: participants, PostStates: postStates}
	tx := Transaction{Data: txdata}

	// Unmarshal
	json.Unmarshal(txbuf, &tx)

	return &tx
}

func (tx *Transaction) Nonce() []uint64 {
	//[FIXME] Get nonce from state? or account nonce field?
	nonces := make([]uint64, 0)
	for _, acc := range tx.Data.PostStates {
		nonces = append(nonces, acc.Nonce)
	}
	return nonces
}

// Get specific user's post state using public key
func (tx *Transaction) GetPostState(key *ecdsa.PublicKey) *state.Account {
        for i, k := range tx.Data.Participants {
                  if *k == *key {
                          return tx.Data.PostStates[i]
                  }
        }
        // No such participant
        return nil
}

func (tx *Transaction) GetPostBalanceSum() uint64 {
          sum := uint64(0)
          for _, s := range tx.Data.PostStates {
                    sum += s.Balance
          }
          return sum
}

func (tx *Transaction) GetPrevBalanceSum() uint64 {
          sum := uint64(0)
          for i, _ := range tx.Data.Participants {
                    //[FIXME]
		    //prevState := loadTransaction(tx.Data.PrevTxHashes[i]).GetPostState()
		    prevState := &state.Account{Balance: uint64(i)}
		    sum += prevState.Balance
          }
          return sum
}

//func (tx *Transaction) Value() uint64 { return tx.Data.Amount }

//func (tx *Transaction) Sender() ecdsa.PublicKey { return *tx.Data.Sender } // Temporal function until signature is implemented

//func (tx *Transaction) Recipient() ecdsa.PublicKey { return *tx.Data.Recipient }

func (tx *Transaction) Participants() []*ecdsa.PublicKey { return tx.Data.Participants }
func (tx *Transaction) PostStates() []*state.Account { return tx.Data.PostStates }
func (tx *Transaction) PrevTxHashes() []*common.Hash { return tx.Data.PrevTxHashes }

// get hashed txdata's byte array
func (data *Txdata) GetHashedBytes() []byte {

	bytelist := []byte{}
	for i := 0; i < len(data.Participants); i++ {
		bytelist = append(bytelist, common.ToBytes(data.Participants[i].X)...)
		bytelist = append(bytelist, common.ToBytes(data.Participants[i].Y)...)
		bytelist = append(bytelist, common.ToBytes(data.PostStates[i].Nonce)...)
		bytelist = append(bytelist, common.ToBytes(data.PostStates[i].Balance)...)
		bytelist = append(bytelist, common.ToBytes(*data.PrevTxHashes[i])...)
	}

	return crypto.Keccak256(bytelist)
}

// hashing txdata of tx
func (tx *Transaction) GetHash() common.Hash {
	//return crypto.Keccak256Hash(common.ToBytes(*tx))

	// new method
	return crypto.Keccak256Hash(tx.Data.GetHashedBytes())
}

// hashing all transactions in txs (temporary)
// TODO: this hash value should be root value of tx's merkle tree
func (txs Transactions) Hash() common.Hash {

	//return crypto.Keccak256Hash(common.ToBytes(*txs))

	// new method
	bytelist := []byte{}
	for i := 0; i < len(txs); i++ {
		bytelist = append(bytelist, txs[i].Data.GetHashedBytes()...)
	}

	return crypto.Keccak256Hash(bytelist)
}

// insert tx into txs
func (txs *Transactions) Insert(tx *Transaction) {
	*txs = append(*txs, tx)
}

func (tx *Transaction) PrintTx() {
	fmt.Println("tx hash:", tx.Hash.ToHex())
	for i := 0; i < len(tx.Data.Participants); i++ {
		fmt.Println("participant ", i)
		fmt.Println("tx hash ", tx.Hash)
		fmt.Println("public key: ", tx.Data.Participants[i])
		//fmt.Println("post state: ", tx.Data.PostStates[i])
		fmt.Print("post state -> ")
		tx.Data.PostStates[i].PrintAccount()
		fmt.Println("previous tx hash: ", tx.Data.PrevTxHashes[i].ToHex())
		fmt.Println()
	}
}

// make random tx for test
func MakeTestTx(participantsNum int, s state.State) *Transaction {
	// make participants
	parNum := participantsNum
	parPrivateKeys := []*ecdsa.PrivateKey{}
	parPublicKeys := []*ecdsa.PublicKey{}
	parStates := []*state.Account{}
	prevTxHashes := []*common.Hash{}

	for i := 0; i < parNum; i++ {
		// make random private/public key pairs
		priv, _ := crypto.GenerateKey()
		parPrivateKeys = append(parPrivateKeys, priv)
		parPublicKeys = append(parPublicKeys, &priv.PublicKey)

		// assume that every participants has 100 ether
		parStates = append(parStates, s.NewAccount(&priv.PublicKey, 0, 100))

		// null prev tx hashes
		prevTxHashes = append(prevTxHashes, &common.Hash{})
	}

	// make transaction
	tx := NewTransaction(parPublicKeys, parStates, prevTxHashes)
	return tx
}

// make random signed tx for test
func MakeTestSignedTx(participantsNum int, s state.State) *Transaction {
	// make participants
	parNum := participantsNum
	parPrivateKeys := []*ecdsa.PrivateKey{}
	parPublicKeys := []*ecdsa.PublicKey{}
	parStates := []*state.Account{}
	prevTxHashes := []*common.Hash{}

	for i := 0; i < parNum; i++ {
		// make random private/public key pairs
		priv, _ := crypto.GenerateKey()
		parPrivateKeys = append(parPrivateKeys, priv)
		parPublicKeys = append(parPublicKeys, &priv.PublicKey)

		// assume that every participants has 100 ether
		parStates = append(parStates, s.NewAccount(&priv.PublicKey, 0, 100))

		// null prev tx hashes
		prevTxHashes = append(prevTxHashes, &common.Hash{})
	}

	// make transaction
	tx := NewTransaction(parPublicKeys, parStates, prevTxHashes)

	// every participants sign to tx
	for i := 0; i < parNum; i++ {
		tx.Sign(parPrivateKeys[i])
	}

	return tx
}

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// tx validation function for iot node
func (tx *Transaction) ValidateTx() error {

	// 1. check Participants, PostStates, PrevTxHashes's lengths are same
	if !(len(tx.Data.Participants) == len(tx.Data.PostStates) && len(tx.Data.PostStates) == len(tx.Data.PrevTxHashes)) {
		return ErrDiffFieldLength
	}

	// 2. check PostStates' Account == Participants' Account (check pub key)
	for i := 0; i < len(tx.Data.Participants); i++ {
		if *tx.Data.PostStates[i].PublicKey != *tx.Data.Participants[i] {
			return ErrInvalidPostStates
		}
	}

	// 3. check PrevTxHashes has Participants state (TODO) -> ErrInvalidPrevTxHashes

	// 4. check signature
	return tx.VerifySignature()
}

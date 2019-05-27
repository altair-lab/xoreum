package rawdb

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/log"
	"github.com/altair-lab/xoreum/xordb"
)

// ReadTxLookupEntry retrieves the blocknumber of the block the tx is in
func ReadTxLookupEntry(db xordb.Reader, hash common.Hash) *uint64 {
	data, _ := db.Get(txLookupKey(hash))
	if len(data) == 0 {
		return nil
	}
	if len(data) < common.HashLength {
		number := new(big.Int).SetBytes(data).Uint64()
		return &number
	}
	return nil
}

// WriteTxLookupEntries stores the blocknumber of the block the tx is in
func WriteTxLookupEntries(db xordb.Writer, block *types.Block) {

	num := common.Num2Bytes(block.Number())
	for _, tx := range block.Transactions() {
		if err := db.Put(txLookupKey(tx.Hash), num); err != nil {
			log.Crit("Failed to store transaction lookup entry", "err", err)
		}
	}
}

// DeleteTxLookupEntry removes all transaction data associated with a hash.
func DeleteTxLookupEntry(db xordb.Writer, hash common.Hash) {
	db.Delete(txLookupKey(hash))
}

// ReadRawTxData retrieves the raw tx data corresponding to the hash
func ReadRawTxData(db xordb.Reader, hash common.Hash) []byte {
	data, _ := db.Get(txRawKey(hash))
	if len(data) == 0 {
		return nil
	}
	if len(data) < common.HashLength {
		return data
	}
	return nil
}

// WriteRawTxData stores the raw tx data corresponding to the hash
func WriteRawTxData(db xordb.Writer, hash common.Hash, data []byte) {
	if err := db.Put(txRawKey(hash), data); err != nil {
		log.Crit("Failed to store raw tx data", "err", err)
	}
}

// DeleteRawTxData removes the raw tx data corresponding to the hash
func DeleteRawTxData(db xordb.Writer, hash common.Hash) {
	db.Delete(txRawKey(hash))
}

// ReadTransaction retrieves a transaction and its metadata.
// returns (tx, blockHash, *blockNumber, uint64(txIndex))
func ReadTransaction(db xordb.Reader, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	// first looks for raw Tx Data.
	rawTxData := ReadRawTxData(db, hash)
	if rawTxData != nil {
		tx := new(types.Transaction)
		json.Unmarshal(rawTxData, &tx)

		txdata := tx.Data
		participants := make([]*(ecdsa.PublicKey), len(txdata.Participants))
		postStates := make([]*(state.Account), len(txdata.Participants))
		for i := 0; i < len(txdata.Participants); i++ {
			txdata.Participants[i] = &ecdsa.PublicKey{Curve: elliptic.P256()}
			txdata.PostStates[i] = &state.Account{PublicKey: txdata.Participants[i]}
		}
		txdata = types.Txdata{Participants: participants, PostStates: postStates}
		json.Unmarshal(rawTxData, &tx)
		return tx, common.Hash{}, 0, 0
	}

	// if none found, looks for the block that contains the tx.
	blockNumber := ReadTxLookupEntry(db, hash)
	if blockNumber == nil {
		return nil, common.Hash{}, 0, 0
	}
	blockHash := ReadHash(db, *blockNumber)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}
	block := LoadBlock(db, blockHash, *blockNumber)
	if block == nil {
		log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash)
		return nil, common.Hash{}, 0, 0
	}
	for txIndex, tx := range block.Transactions() {
		if tx.Hash == hash {
			return tx, blockHash, *blockNumber, uint64(txIndex)
		}
	}
	log.Error("Transaction not found", "number", blockNumber, "hash", blockHash, "txhash", hash)
	return nil, common.Hash{}, 0, 0
}

// WriteTransaction stores a transaction into the database.
func WriteTransaction(db xordb.Writer, hash common.Hash, tx *types.Transaction) {
	data, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("error while encoding", err)
	}
	WriteRawTxData(db, hash, data)
}

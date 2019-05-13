package rawdb

import (
	"fmt"
	"math/big"

	"github.com/altair-lab/xoreum/common"
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

// ReadTransaction retrieves a transaction and its metadata
// (tx, blockHash, *blockNumber, uint64(txIndex))
func ReadTransaction(db xordb.Reader, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {

	blockNumber := ReadTxLookupEntry(db, hash)
	if blockNumber == nil {
		fmt.Println("case1")
		return nil, common.Hash{}, 0, 0
	}
	blockHash := ReadHash(db, *blockNumber)
	if blockHash == (common.Hash{}) {
		fmt.Println("case2")
		return nil, common.Hash{}, 0, 0
	}
	block := LoadBlock(db, blockHash, *blockNumber)
	if block == nil {
		fmt.Println("case3")
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

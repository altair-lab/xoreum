package main

import (

	//"bufio"
	//"strconv"

	"fmt"

	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

const DEFAULT_BLOCK_NUMBER = 10

func main() {

	db, _ := leveldb.New("chaindata", 0, 0, "")

	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	for i := uint64(1); i <= uint64(*last_BN); i++ {
		//load block via accessor api
		fmt.Println("loading block", i)
		loaded := rawdb.LoadBlockByBN(db, i)
		txHash := loaded.Transactions()[0].Hash
		//tx, blockHash, *blockNumber, uint64(txIndex)
		tx, _, blockNumber, txIndex := rawdb.ReadTransaction(db, txHash)
		loaded.PrintBlock()
		fmt.Println("===========")
		if tx == nil {
			fmt.Println("nil ptr")
		} else {
			fmt.Println("tx in block", blockNumber)
			fmt.Println(txIndex, "th transaction")
			// tx.PrintTx()
		}
		fmt.Println("===========")
	}
}

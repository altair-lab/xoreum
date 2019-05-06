package main

import (
	"fmt"

	//"bufio"
	//"strconv"

	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

const DEFAULT_BLOCK_NUMBER = 10

func main() {

	db, _ := leveldb.New("chaindata", 0, 0, "")

	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	for i := uint64(0); i < uint64(*last_BN); i++ {
		//load block via accessor api
		fmt.Println("loading block", i+1)
		rawdb.LoadBlockN(db, i+1)
		fmt.Println("\n")
	}

}

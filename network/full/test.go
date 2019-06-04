/*
  Full Node     : Send all blocks from chain and keep update
  IoT-full Node : Send only interlink blocks from chain and keep update
*/

package main

import (
	"log"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

func main() {
	// Load DB
	db, _ := leveldb.New("db_test", 0, 0, "")

	// Initialize chain and store to DB
	log.Println("Initialize Chain")

	for i := 0; i < 100000; i++ {
		//rawdb.WriteState(db, common.BytesToAddress(common.ToBytes("a")), common.BytesToHash(common.ToBytes("a")))
		rawdb.DeleteState(db, common.BytesToAddress(common.ToBytes("a")))
	}

	log.Println("the number of state : ", rawdb.CountStates(db))
	log.Println("Done")

	db.Close()
}


/*
  Full Node     : Send all blocks from chain and keep update
  IoT-full Node : Send only interlink blocks from chain and keep update
*/

package main

import (
	"log"

	"github.com/altair-lab/xoreum/common"
	//"github.com/altair-lab/xoreum/core/rawdb"
	//"github.com/altair-lab/xoreum/xordb/leveldb"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	// Load DB
	db, _ := leveldb.OpenFile("db_test", nil)
	
	// Initialize chain and store to DB
	log.Println("Initialize Chain")

	for i := 0; i < 100000; i++ {
		//rawdb.WriteState(db, common.BytesToAddress(common.ToBytes("a")), common.BytesToHash(common.ToBytes("a")))
		db.Delete(common.ToBytes("a"), nil)
	}

	//log.Println("the number of state : ", rawdb.CountStates(db))
	log.Println("Done")

	defer db.Close()
}


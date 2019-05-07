/*
  Light Node : Get all blocks from chain and insert
  IoT Node   : Get interlink blocks from chain and validate => Set Genesis block = Currnt block
*/

package main

import (
	"net"
	"log"
	"encoding/json"

	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/network"
)

func main() {
	// create genesis block
	//Blockchain := core.NewBlockChain()

	// Print synchronized json data
	conn, err := net.Dial("tcp","localhost:9000")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}

	for {
		// Make header struct
		buf := network.RecvHeaderJson(conn)
		var header types.Header
		json.Unmarshal([]byte(buf), &header)

		// [TODO] Get Txs json

		// [TODO] Make Txs struct
		txs := types.Transactions{}

		// Make Block with header, txs
		block := types.NewBlock(&header, txs)
		block.Hash() // set block hash
		block.PrintBlock()
		block.PrintTxs()

		// [TODO] State validation (sign, nonce, total balance)
	}
}

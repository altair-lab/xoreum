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
		buf, err := network.RecvObjectJson(conn)
		if err != nil {
			return
		}
		var header types.Header
		json.Unmarshal(buf, &header)
		
		// Get Txs length
		txslen, _ := network.RecvLength(conn)

		// Make Tx struct
		txs := types.Transactions{}
		for i := uint32(0); i < txslen; i++ {
			// Get txdata, R, S
			data, err := network.RecvObjectJson(conn)
			if err != nil {
				return
			}
			R, err := network.RecvObjectJson(conn)
			if err != nil {
				return
			}
			S, err := network.RecvObjectJson(conn)
			if err != nil {
				return
			}

			// Make tx
			tx := types.UnmarshalJSON(data, R, S)
			txs.Insert(tx)
		}

		// Make Block with header, txs
		block := types.NewBlock(&header, txs)
		block.Hash() // set block hash
		block.PrintBlock()
		block.PrintTxs()
		
		// [TODO] State validation (sign, nonce, total balance)
	}
}

/*
  Light Node : Get all blocks from chain and insert
  IoT Node   : Get interlink blocks from chain and validate => Set Genesis block = Currnt block
*/

package main

import (
	"net"
	"log"

	"github.com/altair-lab/xoreum/network"
)

func main() {
	// create genesis block
	// Blockchain := core.NewBlockChain()

	// Print synchronized json data
	conn, err := net.Dial("tcp","localhost:9000")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}

	for {
		block, err := network.RecvBlock(conn)
		if err != nil {
			return
		}

		block.PrintBlock()
		
		// [TODO] State validation (sign, nonce, total balance)
	}
}

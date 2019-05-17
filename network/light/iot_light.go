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

	conn, err := net.Dial("tcp","localhost:9000")
	if nil != err {
		log.Fatal("failed to connect to server")
	}
	
	// Get interlinks length
	interlinkslen, err := network.RecvLength(conn)
	if nil != err {
		log.Fatal(err)
	}

	for i := uint32(0); i < interlinkslen; i++ {
		// Receive interlink block
		block, err := network.RecvBlock(conn)
		if err != nil {
			return
		}
		
		// Block validation (sign, nonce, total balance)
		err = block.ValidateBlock()
		if err != nil{
			log.Fatal(err)
			return
		}

		// Print block
		block.PrintBlock()
	}

	log.Println("INTERLINK SYNCHRONIZATION FINISHED!")

	// [TODO] Make blockchain with current block (= genesis block)
}

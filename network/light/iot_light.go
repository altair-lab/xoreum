/*
  Light Node : Get all blocks from chain and insert
  IoT Node   : Get interlink blocks from chain and validate => Set Genesis block = Currnt block
*/

package main

import (
	"net"
	"log"

	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/xordb/memorydb"
)

var Blockchain *core.BlockChain

func main() {
	// Connect with full node (server)
	conn, err := net.Dial("tcp","localhost:9000")
	if nil != err {
		log.Fatal("failed to connect to server")
	}

	// Receive State
	state, err := network.RecvState(conn)
	if nil != err {
		log.Fatal("failed to receive state")
	}

	// [TODO]
	// Receive Txs (We need temporary 'transactions' object in chain)
	
	// Get interlinks length
	interlinkslen, err := network.RecvLength(conn)
	if nil != err {
		log.Fatal(err)
	}
	currentBlock := &types.Block{} 

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
		currentBlock = block

		// Print block
		block.PrintBlock()
	}

	log.Println("INTERLINK SYNCHRONIZATION FINISHED!")

	// Make IoT blockchain with current block (= genesis block)
	db := memorydb.New()
	Blockchain = core.NewIoTBlockChain(db, currentBlock, state)
	Blockchain.PrintBlockChain()
	Blockchain.GetState().Print()
	Blockchain.GetAllTxs().Print()
/*
	// [TODO] Keep mining every MINING_INTERVAL
	go func() {
		for {
			time.Sleep(MINING_INTERVAL * time.Second)
			// Mining from txpool
			block := Miner.Mine(Txpool, uint64(0))
			if block != nil {
				block.PrintTxs()
			} else {
				fmt.Println("Mining Fail")
			}
			// Add to Blockchain
			err := Blockchain.Insert(block)
			if err != nil {
				fmt.Println(err)
			}
			Blockchain.CurrentBlock().PrintBlock()
		}
	}()
*/
}

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
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

var Blockchain *core.BlockChain

func main() {
	// Load DB
	db, _ := leveldb.New("iot-chaindata", 0, 0, "")
	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	// When there is no existing DB
	if last_BN == nil {
		// Connect with full node (server)
		conn, err := net.Dial("tcp","localhost:9000")
		if nil != err {
			log.Fatal("failed to connect to server")
		}

		// Receive State
		state, allTxs, err := network.RecvState(conn)
		if nil != err {
			log.Fatal("failed to receive state")
		}
	
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

		// Make IoT blockchain with current block (= genesis block)
		Blockchain = core.NewIoTBlockChain(db, currentBlock, state, allTxs)
	} else {
		// Load blocks via accessor api
		// [TODO] Get Genesis block by GenesisBN
		genesis := rawdb.LoadBlockByBN(db, *last_BN)
		// [TODO] Fix after removing state, allTxs fields
		Blockchain = core.NewIoTBlockChain(db, genesis, nil, nil)
		/*
		for i := Blockchain.Genesis().GetHeader().Number; i <= *last_BN; i++ {
			loaded := rawdb.LoadBlockByBN(db, i)
			err := Blockchain.Insert(loaded)
			if err != nil {
				log.Println(err)
				return
			}
		}
		*/
	}

	// Print blockchain
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

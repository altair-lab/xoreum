/*
  Light Node : Get all blocks from chain and insert
  IoT Node   : Get interlink blocks from chain and validate => Set Genesis block = Currnt block
*/

package main

import (
	"log"
	"net"
	"os"
	"time"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

var Blockchain *core.BlockChain

func main() {


	// Load DB
	db, _ := leveldb.New("chaindata-iot", 0, 0, "")
	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	// When there is no existing DB
	if last_BN == nil {
		// Connect with full node (server)
		host := "lynx.snu.ac.kr"
		port := "8085" // Default port number
		if len(os.Args) > 1 {
			port = os.Args[1]
		}

		start := time.Now()
		conn, err := net.Dial("tcp", host+":"+port)
		if nil != err {
			log.Fatal("failed to connect to server")
		}
		
		log.Println("Connected!")

		elapsed := time.Since(start)
		log.Printf("%s", elapsed)
		// Receive State
		err = network.RecvState(conn, db)
		if nil != err {
			log.Fatal("failed to receive state")
		}
		log.Println("Receive state done!")
		
		elapsed = time.Since(start)
		log.Printf("%s", elapsed)
		
		// Get interlinks length
		interlinkslen, err := network.RecvLength(conn)
		if nil != err {
			log.Fatal(err)
		}
		currentBlock := &types.Block{}
		
		log.Println("Receive Interlink Blocks . . .")
		for i := uint32(0); i < interlinkslen; i++ {
			// Receive interlink block
			block, err := network.RecvBlock(conn)
			if err != nil {
				return
			}

			// Block validation (sign, nonce, total balance)
			err = block.ValidateBlock()
			if err != nil {
				log.Fatal(err)
				return
			}
			currentBlock = block

			// Print block
			//block.PrintBlock()
		}

		elapsed = time.Since(start)
		log.Printf("%s", elapsed)
		// Make IoT blockchain with current block (= genesis block)
		Blockchain = core.NewIoTBlockChain(db, currentBlock)
		rawdb.WriteLastHeaderHash(db, currentBlock.GetHeader().Hash())
		log.Println("Synchronization Done!")

		elapsed = time.Since(start)
		log.Printf("%s", elapsed)


	} else {
		// Load blocks after genesis block
		genesis_hash := rawdb.ReadGenesisHeaderHash(db)
		genesis_BN := rawdb.ReadHeaderNumber(db, genesis_hash)
		genesis := rawdb.LoadBlockByBN(db, *genesis_BN)
		Blockchain = core.NewIoTBlockChain(db, genesis)
		log.Println("Load Block Done!")
	}

	// Print blockchain
	//Blockchain.PrintBlockChain()
	//rawdb.ReadStates(db)


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

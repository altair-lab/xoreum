/*
  Light Node : Get all blocks from chain and insert
  IoT Node   : Get interlink blocks from chain and validate => Set Genesis block = Currnt block
*/

package main

import (
	"log"
	"net"
	"os"
	"encoding/json"
	"path/filepath"

	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

var Blockchain *core.BlockChain

type Configuration struct {
	Hostname	string
	Port		string
	BlockNumber	int64
	Participants	int64
	PrintMode	bool
	MiningInterval	int
}

func main() {
	// Load configuration
        ex, err := os.Executable()
        if err != nil {
                log.Println(err)
        }
        exPath := filepath.Dir(ex)
        file, _ := os.Open(exPath+"/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Println("error : ", err)
	}

	// Load DB
	db, _ := leveldb.New("chaindata-iot", 0, 0, "")
	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	// When there is no existing DB
	if last_BN == nil {
		// Connect with full node (server)
		conn, err := net.Dial("tcp", configuration.Hostname+":"+configuration.Port)
		if nil != err {
			log.Fatal("failed to connect to server")
		}

		log.Println("Conntected!")

		// Receive State
		err = network.RecvState(conn, db)
		if nil != err {
			log.Fatal("failed to receive state")
		}
		log.Println("Receive state done!")

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

			// Print received blocks
			if (configuration.PrintMode) {
				block.PrintBlock()
			}
		}

		// Make IoT blockchain with current block (= genesis block)
		Blockchain = core.NewIoTBlockChain(db, currentBlock)
		rawdb.WriteLastHeaderHash(db, currentBlock.GetHeader().Hash())
		log.Println("Synchronization Done!")
	} else {
		// Load blocks after genesis block
		genesis_hash := rawdb.ReadGenesisHeaderHash(db)
		genesis_BN := rawdb.ReadHeaderNumber(db, genesis_hash)
		genesis := rawdb.LoadBlockByBN(db, *genesis_BN)
		Blockchain = core.NewIoTBlockChain(db, genesis)
		log.Println("Load Block Done!")
	}

	// Print blockchain
	if (configuration.PrintMode) {
		Blockchain.PrintBlockChain()
		//rawdb.ReadStates(db)
	}

	// [TODO] Keep mining every interval
}

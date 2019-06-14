/*
  IoT-full Node : Send only interlink blocks from chain and keep update
*/

package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
	"encoding/json"

	"github.com/altair-lab/xoreum/xordb"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)
/*
const MINING_INTERVAL = 10
const DEFAULT_BLOCK_NUMBER = 255
const DEFAULT_ACCOUNTS_NUMBER = 64
const BROADCAST_INTERVAL = 5
*/
var Blockchain *core.BlockChain
var mutex = &sync.Mutex{}

type Configuration struct {
	Hostname	string
	Port		string
	BlockNumber	int64
	Participants	int64
	Difficulty	int
	MiningInterval	int
	BroadcastInterval	int
}

func main() {
	// Load configuration
	file, _ := os.Open("../conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("error : ", err)
	}
	log.Println(configuration)

	// Load DB
	db, _ := leveldb.New("chaindata", 0, 0, "")
	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	// When there is no existing DB
	if last_BN == nil {
		// Initialize chain and store to DB
		log.Println("Initialize Chain")
		Blockchain = network.MakeTestBlockChain(configuration.BlockNumber, configuration.Participants, db)
		log.Println("Done")
	} else {
		// Load blocks from 1st block (0 = genesis)
		log.Println("Load Chain")
		Blockchain = core.NewBlockChain(db)
		log.Println("Done")
	}

	// Print blckchain
	Blockchain.PrintBlockChain()
	//rawdb.ReadStates(db)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", configuration.Hostname+":"+configuration.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	// Create a new connection each time we receive a connection request
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn, db)
	}
}

// connection
func handleConn(conn net.Conn, db xordb.Database) {
	addr := conn.RemoteAddr().String()
	recvBuf := make([]byte, 4096) // receive buffer: 4kb

	// Connected to new client
	log.Printf("CONNECTED TO %v\n", addr)

	// Remove GetAccounts
	network.SendState(conn, db)

	// Send only Interlink block data
	interlinks := Blockchain.CurrentBlock().GetUniqueInterlink()
	network.SendInterlinks(conn, interlinks, Blockchain)
	quit := make(chan bool)

	// Check recvBuf every clock
	go func() {
		for {
			// Get input data from clients every clock
			_, err := conn.Read(recvBuf)

			if nil != err {
				if io.EOF == err {
					log.Printf("Connection is closed from client; %v", conn.RemoteAddr().String())
					quit <- true
					return
				}
				log.Printf("fail to receive data; err: %v", err)
				quit <- true
				return
			}
		}

	}()
}

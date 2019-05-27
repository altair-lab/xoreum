/*
  Full Node     : Send all blocks from chain and keep update
  IoT-full Node : Send only interlink blocks from chain and keep update
*/

package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/altair-lab/xoreum/xordb"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

const MINING_INTERVAL = 10
const DEFAULT_BLOCK_NUMBER = 5
const DEFAULT_ACCOUNTS_NUMBER = 5
const BROADCAST_INTERVAL = 5

var Blockchain *core.BlockChain
var mutex = &sync.Mutex{}

func main() {
	// Load DB
	db, _ := leveldb.New("chaindata", 0, 0, "")
	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	// When there is no existing DB
	if last_BN == nil {
		// Initialize chain and store to DB
		Blockchain = network.MakeTestBlockChain(DEFAULT_BLOCK_NUMBER, DEFAULT_ACCOUNTS_NUMBER, db)
		log.Println("Initialize Chain")
	} else {
		// Load blocks from 1st block (0 = genesis)
		Blockchain = core.NewBlockChain(db)
		for i := uint64(1); i <= *last_BN; i++ {
			loaded := rawdb.LoadBlockByBN(db, i)
			err := Blockchain.Insert(loaded)
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println("Load Chain")
	}

	// Print blckchain
	Blockchain.PrintBlockChain()
	//Blockchain.GetAccounts().Print()
	//Blockchain.GetAllTxs().Print()

	// start TCP and serve TCP server
	port := "9000" //  Default port number
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	server, err := net.Listen("tcp", ":"+port)
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

	// [FIXME] Send State
	// Remove GetAccounts
	network.SendState(conn, db, Blockchain.GetAccounts())

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

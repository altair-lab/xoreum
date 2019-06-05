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
	"strconv"

	"github.com/altair-lab/xoreum/xordb"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/xordb/leveldb"
)

const MINING_INTERVAL = 10
const DEFAULT_BLOCK_NUMBER = 1000
const DEFAULT_ACCOUNTS_NUMBER = 10000
const BROADCAST_INTERVAL = 5

var Blockchain *core.BlockChain
var mutex = &sync.Mutex{}

func main() {
	// Set Default Block Number for test
    	testBlockNum := int64(DEFAULT_BLOCK_NUMBER)
    	if len(os.Args) > 1 {
        	testBlockNum, _ = strconv.ParseInt(os.Args[1], 10, 64)
	}

	// Set Default Account Number for test
    	testAccNum := int64(DEFAULT_ACCOUNTS_NUMBER)
    	if len(os.Args) > 2 {
        	testAccNum, _ = strconv.ParseInt(os.Args[2], 10, 64)
	}

	// Load DB
	db, _ := leveldb.New("chaindata_"+strconv.FormatInt(testBlockNum, 10), 0, 0, "")
	last_hash := rawdb.ReadLastHeaderHash(db)
	last_BN := rawdb.ReadHeaderNumber(db, last_hash)

	// When there is no existing DB
	if last_BN == nil {
		// Initialize chain and store to DB
		log.Println("Initialize Chain")
		Blockchain = network.MakeTestBlockChain(testBlockNum, testAccNum, db)
		log.Println("Done")
	} else {
		// Load blocks from 1st block (0 = genesis)
		log.Println("Load Chain")
		Blockchain = core.NewBlockChain(db)
		log.Println("#blocks: ", *last_BN)
		log.Println("#states: ", rawdb.CountStates(db))
		log.Println("Done")
		return
	}

	// Print blckchain
	//Blockchain.PrintBlockChain()
	//rawdb.ReadStates(db)

	// start TCP and serve TCP server
	host := ""
	port := "8084" //  Default port number (yj:8084, yh:8085)
/*
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
*/
	server, err := net.Listen("tcp", host+":"+port)
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

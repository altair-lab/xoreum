/*
  Full Node     : Send all blocks from chain and keep update
  IoT-full Node : Send only interlink blocks from chain and keep update
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/miner"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/xordb/memorydb"
)

const MINING_INTERVAL = 10
const DEFAULT_BLOCK_NUMBER = 10
const BROADCAST_INTERVAL = 5

// [TODO] replaceChain (logest chain rule)
var Blockchain *core.BlockChain
var Txpool *core.TxPool
var Miner miner.Miner
var mutex = &sync.Mutex{}

func main() {
	// create genesis block
	db := memorydb.New()

	Blockchain = core.NewBlockChain(db)
	Blockchain.PrintBlockChain()

	// Initialization txpool, miner
	Txpool, Miner = network.Initialization(Blockchain)

	// Mining default blocks (for test)
	for i := 0; i < DEFAULT_BLOCK_NUMBER; i++ {
		// Make test tx and add to txpool
		for j := 0; j < i; j++ {
			success, err := Txpool.Add(types.MakeTestSignedTx(j + 1, Blockchain.GetState()))
			if !success {
				fmt.Println(err)
			}
		}

		// Make block (mining)
		block := Miner.Mine(Txpool, uint64(0))
		if block == nil {
			fmt.Println("Mining Fail")
		}

		// Insert block to Blockchain
		err := Blockchain.Insert(block)
		if err != nil {
			fmt.Println(err)
		}

		// Print current block
		Blockchain.CurrentBlock().PrintBlock()
		Blockchain.CurrentBlock().PrintTxs()
	}

	// Keep mining every MINING_INTERVAL
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

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":9000")
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
		go handleConn(conn)
	}
}

// connection
func handleConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	recvBuf := make([]byte, 4096) // receive buffer: 4kb

	// Connected to new client
	log.Printf("CONNECTED TO %v\n", addr)

	// Send only Interlink block data
	currentBlockNumber := Blockchain.CurrentBlock().GetHeader().Number
	updatedBlockNumber := uint64(0)
	interlinks := Blockchain.CurrentBlock().GetUniqueInterlink()
	log.Printf("INTERLINKS : %v\n", interlinks)
	for i := 0; i < len(interlinks); i++ {
		// Send block
		err := network.SendBlock(conn, Blockchain.BlockAt(interlinks[i]))
		if err != nil {
			return
		}
		updatedBlockNumber = interlinks[i]
	}

	quit := make(chan bool)

	// Check recvBuf every clock
	go func() {
		for {
			// Get input data from clients every clock
			n, err := conn.Read(recvBuf)

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

			if 0 < n {
				data := recvBuf[:n]
				log.Println(string(data))
				// [TODO] Make Tx ?
			}
		}

	}()

	// Send newly updated block (check every BROADCAST_INTERVAL)
	go func() {
		for {
			time.Sleep(BROADCAST_INTERVAL * time.Second)

			select {
			case <-quit:
				return
			default:
				// Check updated block
				currentBlockNumber = Blockchain.CurrentBlock().GetHeader().Number
				for i := updatedBlockNumber + 1; i <= Blockchain.CurrentBlock().GetHeader().Number; i++ {
					// Send block
					err := network.SendBlock(conn, Blockchain.BlockAt(i))
					if err != nil {
						return
					}
					updatedBlockNumber = i
				}
			}
		}
	}()
}

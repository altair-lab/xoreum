/*
  Full Node     : Send all blocks from chain and keep update
  IoT-full Node : Send only interlink blocks from chain and keep update
*/

package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/network"
)

const MINING_INTERVAL = 10
const DEFAULT_BLOCK_NUMBER = 10
const BROADCAST_INTERVAL = 5

// [TODO] replaceChain (logest chain rule)
var Blockchain *core.BlockChain
var mutex = &sync.Mutex{}

func main() {
	// create genesis block
	//db := memorydb.New()
	Blockchain = network.MakeTestBlockChain(3, 5)
	Blockchain.PrintBlockChain()
/*
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
*/
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
	interlinks := Blockchain.CurrentBlock().GetUniqueInterlink()
	network.SendInterlinks(conn, interlinks, Blockchain)
	updatedBlockNumber := interlinks[len(interlinks)-1]
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

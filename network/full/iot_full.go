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

	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/network"
)

const MINING_INTERVAL = 10
const DEFAULT_BLOCK_NUMBER = 5
const DEFAULT_ACCOUNTS_NUMBER = 5
const BROADCAST_INTERVAL = 5

// [TODO] replaceChain (logest chain rule)
var Blockchain *core.BlockChain
var mutex = &sync.Mutex{}

func main() {
	// create block chain
	Blockchain = network.MakeTestBlockChain(DEFAULT_BLOCK_NUMBER, DEFAULT_ACCOUNTS_NUMBER)
	Blockchain.PrintBlockChain()

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
	Blockchain.GetState().Print()
	Blockchain.GetAllTxs().Print()

	// Send State
	network.SendState(conn, Blockchain.GetState(), Blockchain.GetAllTxs())

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

/*
  Full Node     : Send all blocks from chain and keep update
  IoT-full Node : Send only interlink blocks from chain and keep update
*/

package main 

import (
	"fmt"
	"log"
	"net"
	"time"
	"io"
	"sync"
	//"bufio"
	//"strconv"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/network"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/miner"
)

const MINING_INTERVAL = 10
const DEFAULT_BLOCK_NUMBER = 10
const BROADCAST_INTERVAL = 5

// [TODO] replaceChain (logest chain rule)
var Blockchain *core.BlockChain
var Acc0 *state.Account
var State state.State
var Txpool *core.TxPool
var Miner miner.Miner

// bcServer handles incoming concurrent Blocks
var bcServer chan *core.BlockChain
var mutex = &sync.Mutex{}


func main() {
	bcServer = make(chan *core.BlockChain)

	// create genesis block
	Blockchain = core.NewBlockChain()
	Blockchain.PrintBlockChain()
	
	// set account, txpool, state, miner for mining
	privatekey0, _ := crypto.GenerateKey()
	publickey0 := privatekey0.PublicKey
	address0 := crypto.Keccak256Address(common.ToBytes(publickey0))
	Acc0 := state.NewAccount(address0, uint64(0), uint64(7000)) // acc0 [Address:0, Nonce:0, Balance:7000]
	State = state.NewState()
	State.Add(Acc0)
	Txpool = core.NewTxPool(State, Blockchain)
	Miner = miner.Miner{Acc0.Address}

	for i := 0; i < DEFAULT_BLOCK_NUMBER; i++ {
		// Make test tx and add to txpool
		success, err := Txpool.Add(types.MakeTestSignedTx(2))
		if !success {
			fmt.Println(err)
		}

		// Make block (mining)
		block := Miner.Mine(Txpool, uint64(0))
		if block == nil {
       			fmt.Println("Mining Fail")
       		}

      		// Insert block to Blockchain
		err = Blockchain.Insert(block)
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

		// Send block header
		network.SendObject(conn, Blockchain.BlockAt(interlinks[i]).GetHeader())
		updatedBlockNumber = interlinks[i]

		// [TODO] Send transactions txdata
		// SendObject(conn, Blockchain.BlockAt(interlinks[i]).GetTxs())
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
				case <- quit:
					return
				default:
				// Check updated block
				currentBlockNumber = Blockchain.CurrentBlock().GetHeader().Number
				for i := updatedBlockNumber + 1; i <= Blockchain.CurrentBlock().GetHeader().Number; i++ {
					network.SendObject(conn, Blockchain.BlockAt(i).GetHeader())
					updatedBlockNumber = i
				}
			}
		}
	}()
}

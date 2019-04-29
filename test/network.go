package main 

import (
	"fmt"
	"encoding/json"
	"log"
	"net"
	"time"
	"io"
	"sync"
	"bufio"
	"strconv"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/miner"

	"github.com/joho/godotenv"
)

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan *core.BlockChain)

	// create genesis block
	Blockchain = core.NewBlockChain()
	
	// set account, txpool, state, miner for mining
	privatekey0, _ := crypto.GenerateKey()
	publickey0 := privatekey0.PublicKey
	address0 := crypto.Keccak256Address(common.ToBytes(publickey0))
	Acc0 := state.NewAccount(address0, uint64(0), uint64(7000)) // acc0 [Address:0, Nonce:0, Balance:7000]
	State = state.NewState()
	State.Add(Acc0)
	Txpool = core.NewTxPool(State, Blockchain)
	Miner = miner.Miner{Acc0.Address}

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

func handleConn(conn net.Conn) {
	defer conn.Close()
	
	// [TODO] Create transaction using this
	//        Not block (block will be created periodically)
	io.WriteString(conn, "Enter a difficulty: ")

	scanner := bufio.NewScanner(conn)

	// Add it to blockchain after conducting validation
	go func() {
		for scanner.Scan() {
			// Mining from txpool
			inputNum, _ := strconv.Atoi(scanner.Text())
			block := Miner.Mine(Txpool, uint64(inputNum))

			if block != nil {
       	        	 	block.PrintTx()
        		} else {
                		fmt.Println("Mining Fail")
        		}

        		// Add to Blockchain
			err := Blockchain.Insert(block)
        		if err != nil {
                		fmt.Println(err)
        		}

			bcServer <- Blockchain
		}
	}()

	// Simulate receiving broadcast
	go func() {
		for {
			// client output
			time.Sleep(5 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain.CurrentBlock().GetHeader())
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output)+"\n")
		}
	}()

	for _ = range bcServer {
		Blockchain.PrintBlockChain()
	}
}


package main 

import (
	"fmt"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
	"io"
	"sync"
	"bufio"

	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/joho/godotenv"
)

// [TODO] replaceChain (logest chain rule)

var Blockchain *core.BlockChain

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

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
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

	// output:
	// true

}

func handleConn(conn net.Conn) {
	// [TODO] Making block by mining

	/*
	privatekey0, _ := crypto.GenerateKey()
	publickey0 := privatekey0.PublicKey
	address0 := crypto.Keccak256Address(common.ToBytes(publickey0))
	acc0 := state.NewAccount(address0, uint64(0), uint64(7000)) // acc0 [Address:0, Nonce:0, Balance:7000]
	state := state.NewState()
	state.Add(acc0)
	*/	
	defer conn.Close()

	io.WriteString(conn, "Enter a value:")

	scanner := bufio.NewScanner(conn)

	// Make Tx and Mining
	// Add it to blockchain after conducting validation
	go func() {
		for scanner.Scan() {
			empty_txs := []*types.Transaction{}
			block := types.NewBlock(&types.Header{}, empty_txs)
			block.GetHeader().ParentHash = Blockchain.CurrentBlock().Hash()
			block.GetHeader().Number = Blockchain.CurrentBlock().GetHeader().Number+1
			block.GetHeader().Nonce = 0
			block.GetHeader().InterLink = Blockchain.CurrentBlock().GetUpdatedInterlink()
			block.GetHeader().Difficulty = 1

			/*
			// Mining from txpool
        		miner := miner.Miner{acc0.Address}
        		block := miner.Mine(txpool, state, uint64(scanner.Text()))
			*/

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
			time.Sleep(30 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcServer {
		Blockchain.PrintBlockChain()
	}
}



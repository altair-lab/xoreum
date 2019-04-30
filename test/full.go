package main 

import (
	"fmt"
	"encoding/json"
	"log"
	"net"
	"time"
	"io"
	"sync"
	//"bufio"
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

// connection : go run network_client.go
func handleConn(conn net.Conn) {
	recvBuf := make([]byte, 4096) // receive buffer: 4kb
	
	go func() {
		for {
			n, err := conn.Read(recvBuf)

			if nil != err {
				if io.EOF == err {
					log.Printf("Connection is closed from client; %v", conn.RemoteAddr().String())
					return
				}
				log.Printf("fail to receive data; err: %v", err)
				return
			}

			if 0 < n {
				data := recvBuf[:n]
				log.Println(string(data))

				// Mining from txpool
				inputNum, _ := strconv.Atoi(string(data))
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

				Blockchain.PrintBlockChain()
				bcServer <- Blockchain
			}
		}

	}()

	go func() {
		for {
			// client output
			time.Sleep(5 * time.Second)
			for i := uint64(0); i <= Blockchain.CurrentBlock().GetHeader().Number; i++ {
				mutex.Lock()
				output, err := json.Marshal(Blockchain.BlockAt(i).GetHeader())
				if err != nil {
					log.Fatal(err)
				}
				mutex.Unlock()
				conn.Write([]byte(string(output)+"\n"))
			}
		}
	}()
}

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	//"bufio"
	//"strconv"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/miner"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/xordb/leveldb"
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

	db, _ := leveldb.New("chaindata", 0, 0, "")

	// create genesis block
	Blockchain = core.NewBlockChain(db)

	// set account, txpool, state, miner for mining
	privatekey0, _ := crypto.GenerateKey()
	publickey0 := privatekey0.PublicKey
	address0 := crypto.Keccak256Address(common.ToBytes(publickey0))
	Acc0 := state.NewAccount(address0, uint64(0), uint64(7000)) // acc0 [Address:0, Nonce:0, Balance:7000]
	State = state.NewState()
	State.Add(Acc0)
	Txpool = core.NewTxPool(State, Blockchain)
	Miner = miner.Miner{Acc0.Address}
	last_BN := uint64(0)
	for i := uint64(1); i < uint64(DEFAULT_BLOCK_NUMBER+1); i++ {
		Txpool.Add(types.MakeTestSignedTx(2))

		block := Miner.Mine(Txpool, uint64(0))

		if block != nil {
			// block.PrintTxs()
		} else {
			fmt.Println("Mining Fail")
		}

		// Add to Blockchain
		err := Blockchain.Insert(block)
		if err != nil {
			fmt.Println(err)
		}

		//store block via rawdb accessor api
		fmt.Println("storing block", block.Number())
		rawdb.StoreBlock(db, block)
		block.PrintBlock()
		fmt.Println("\n")

		last_BN = i
		last_hash := rawdb.ReadHash(db, last_BN)
		rawdb.WriteLastHeaderHash(db, last_hash)
		fmt.Println("last:", last_BN)
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

			fmt.Println("storing block", block.Number())
			rawdb.StoreBlock(db, block)
			fmt.Println("\n")

			last_BN = block.Number()
			last_hash := rawdb.ReadHash(db, last_BN)
			rawdb.WriteLastHeaderHash(db, last_hash)
			fmt.Println("last:", last_BN)
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

// Send message with size
func SendMessage(conn net.Conn, msg []byte) error {
	lengthBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBuf, uint32(len(msg)))
	if _, err := conn.Write(lengthBuf); nil != err {
		log.Printf("failed to send msg length; err: %v", err)
		return err
	}

	if _, err := conn.Write(msg); nil != err {
		log.Printf("failed to send msg; err: %v", err)
		return err
	}

	return nil
}

// connection
func handleConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	recvBuf := make([]byte, 4096) // receive buffer: 4kb

	// Connected to new client
	log.Printf("CONNECTED TO %v\n", addr)

	// Send full block data once
	currentBlockNumber := Blockchain.CurrentBlock().GetHeader().Number
	updatedBlockNumber := uint64(0)
	for i := uint64(0); i <= currentBlockNumber; i++ {
		mutex.Lock()
		output, err := json.Marshal(Blockchain.BlockAt(i).GetHeader())
		if err != nil {
			log.Fatal(err)
		}
		mutex.Unlock()
		log.Printf("Block Length : %d\n", len(output))
		err = SendMessage(conn, output)
		if err != nil {
			log.Fatal(err)
		}
		updatedBlockNumber = i
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
					mutex.Lock()
					output, err := json.Marshal(Blockchain.BlockAt(i).GetHeader())
					if err != nil {
						log.Fatal(err)
					}
					mutex.Unlock()
					log.Printf("Block Length : %d\n", len(output))
					err = SendMessage(conn, output)
					if err != nil {
						log.Fatal(err)
					}
					updatedBlockNumber = i
				}
			}
		}
	}()
}

/*
  Light Node : Get all blocks from chain and insert
  IoT Node   : Get interlink blocks from chain and validate => Set Genesis block = Currnt block
*/

package main

import (
	//"fmt"
	"io"
	"log"
	"net"

	//"strings"
	"encoding/binary"
	"encoding/json"

	//"time"

	"github.com/altair-lab/xoreum/core/types"
	//"github.com/altair-lab/xoreum/core"
)

func main() {
	// create genesis block
	//Blockchain := core.NewBlockChain()

	// Print synchronized json data
	conn, err := net.Dial("tcp", "localhost:9000")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}

	for {
		// Get header json
		length, err := RecvLength(conn)
		if err != nil {
			if io.EOF == err {
				log.Printf("Connection is closed from server; %v", conn.RemoteAddr().String())
				return
			}
			log.Fatal(err)
		}
		buf := make([]byte, length)
		_, err = conn.Read(buf)
		if err != nil {
			if io.EOF == err {
				log.Printf("Connection is closed from server; %v", conn.RemoteAddr().String())
				return
			}
			log.Fatal(err)
		}

		// Make header struct
		var header types.Header
		json.Unmarshal([]byte(buf), &header)

		// [TODO] Get Txs json

		// [TODO] Make Txs struct
		txs := types.Transactions{}

		// Make Block with header, txs
		block := types.NewBlock(&header, txs)
		block.Hash() // set block hash
		block.PrintBlock()

		// [TODO] State validation (sign, nonce, total balance)
	}
}

func RecvLength(conn net.Conn) (uint32, error) {
	lengthBuf := make([]byte, 4)
	_, err := conn.Read(lengthBuf)
	if nil != err {
		return uint32(0), err
	}

	msgLength := binary.LittleEndian.Uint32(lengthBuf)

	return uint32(msgLength), err
}

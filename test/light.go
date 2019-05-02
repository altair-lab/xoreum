package main

import (
	"fmt"
	"net"
	"log"
	//"strings"
	"encoding/binary"
	"encoding/json"
	//"time"

	"github.com/altair-lab/xoreum/core/types"
)

func main() {
	// Print synchronized json data
	conn, err := net.Dial("tcp","localhost:9000")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}

	for {
		// Get header json
		length, err := RecvLength(conn)
		log.Println("length: ", length)
		buf := make([]byte, length)
		_, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(buf))

		// Make header struct
		var h types.Header
		json.Unmarshal([]byte(buf), &h)
		fmt.Println(h)

		// [TODO] Get Txs json

		// [TODO] Make Txs struct

		// [TODO] Check server connection
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

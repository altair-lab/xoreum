package main

import (
	"fmt"
	"net"
	"log"
	//"time"
)

func main() {
	// Print synchronized json data
	conn, err := net.Dial("tcp","localhost:9000")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}

	for {
		// [TODO] ISSUE : buffer size to get block data
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(buf))
	}
}

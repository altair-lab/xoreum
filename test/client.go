package main

import (
	"net"
	"log"
	"bufio"
	"os"
	"fmt"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		conn, err := net.Dial("tcp","localhost:9000")
		if nil != err {
			log.Fatalf("failed to connect to server")
		}
		
		fmt.Print("Enter Difficulty: ")
		value, _ := reader.ReadString('\n')
		value = strings.TrimSuffix(value, "\n")
		msg := []byte(value)
		_, err = conn.Write(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

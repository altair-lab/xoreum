package network

import (
	"net"
	"encoding/binary"
	"log"
	"encoding/json"
	"sync"
)

var mutex = &sync.Mutex{}

type object interface {
//	Unmarshal()
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

// Send object with handling mutex, err
func SendObject(conn net.Conn, v interface{}) {
        mutex.Lock()
        output, err := json.Marshal(v)
        if err != nil {
                log.Fatal(err)
                return
        }
        mutex.Unlock()
        err = SendMessage(conn, output)
        if err != nil {
                log.Fatal(err)
                return
        }
}

// Receive message size
func RecvLength(conn net.Conn) (uint32, error) {
        lengthBuf := make([]byte, 4)
        _, err := conn.Read(lengthBuf)
        if nil != err {
                return uint32(0), err
        }

        msgLength := binary.LittleEndian.Uint32(lengthBuf)

        return uint32(msgLength), err
}

// Unmarshal object message


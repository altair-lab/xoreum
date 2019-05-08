package network

import (
	"net"
	"encoding/binary"
	"log"
	"encoding/json"
	"sync"
	"io"

	"github.com/altair-lab/xoreum/core/types"
)

var mutex = &sync.Mutex{}

// Send message with buffer size
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
func SendObject(conn net.Conn, v interface{}) error {
        mutex.Lock()
        output, err := json.Marshal(v)
        if err != nil {
                log.Fatal(err)
                return err
        }
        mutex.Unlock()
        err = SendMessage(conn, output)
        if err != nil {
                log.Fatal(err)
                return err
        }
	return nil
}

// [TODO] Send Transaction with Signature and txdata
func SendTransactions(conn net.Conn, txs *types.Transactions) error {
	// Send txs length
	txslen := make([]byte, 4)
	binary.LittleEndian.PutUint32(txslen, uint32(len(*txs)))
	if _, err := conn.Write(txslen); nil != err {
		log.Printf("failed to send tx length; err: %v", err)
		return err
	}

	// Send txs
	for i := 0; i < len(*txs); i++ {
		// 1. Send txdata
		mutex.Lock()
		output, err:= json.Marshal((*txs)[i].Data())
		if err != nil {
			log.Fatal(err)
			return err
		}
		mutex.Unlock()
		err = SendMessage(conn, output)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// 2. Send Signatures (R)
		mutex.Lock()
		output, err = json.Marshal((*txs)[i].Signature_R)
		if err != nil {
			log.Fatal(err)
			return err
		}
		mutex.Unlock()
		err = SendMessage(conn, output)
		if err != nil {
			log.Fatal(err)
			return err
		}
		
		// 3. Send Signature (S)
		mutex.Lock()
		output, err = json.Marshal((*txs)[i].Signature_S)
		if err != nil {
			log.Fatal(err)
			return err
		}
		mutex.Unlock()
		err = SendMessage(conn, output)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	
	return nil
}

// Send Block
func SendBlock(conn net.Conn, block *types.Block) error {
	header := block.GetHeader()
	txs := block.GetTxs()
	
	err := SendObject(conn, header)
	if err != nil {
		return err
	}
	err = SendTransactions(conn, txs)
	if err != nil {
		return err
	}

	return nil
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

// Get object json
func RecvObjectJson(conn net.Conn) ([]byte, error) {
	length, err := RecvLength(conn)
        if err != nil {
        	if io.EOF == err {
                	log.Printf("Connection is closed from server; %v", conn.RemoteAddr().String())
                        return nil, err
                }
                log.Fatal(err)
		return nil, err
        }
        buf := make([]byte, length)
        _, err = conn.Read(buf)
        if err != nil {
        	if io.EOF == err {
                	log.Printf("Connection is closed from server; %v", conn.RemoteAddr().String())
                        return nil, err
                }
                log.Fatal(err)
		return nil, err
        }
	return buf, nil
}

// Get Transaction object json
func RecvBlock(conn net.Conn) (*types.Block, error) {
	// Make header struct
	buf, err := RecvObjectJson(conn)
	if err != nil {
		return nil, err
	}
	var header types.Header
	json.Unmarshal(buf, &header)

	// Get Txs length
	txslen, err := RecvLength(conn)
	if err != nil {
		return nil, err
	}
	
	// Make Tx struct
	txs := types.Transactions{}
	for i := uint32(0); i < txslen; i++ {
		// Get txdata, R, S
		data, err := RecvObjectJson(conn)
		if err != nil {
			return nil, err
		}
		
		R, err := RecvObjectJson(conn)
		if err != nil {
			return nil, err
		}
		
		S, err := RecvObjectJson(conn)
		if err != nil {
			return nil, err
		}

		// Unmarshal Tx
		tx := types.UnmarshalJSON(data, R, S)
		txs.Insert(tx)
	}

	// Make Block with header, txs
	block := types.NewBlock(&header, txs)
	block.Hash() // set block hash
	return block, nil
}

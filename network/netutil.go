package network

import (
	"runtime"
	"fmt"
	"net"
	"encoding/binary"
	"log"
	"encoding/json"
	"sync"
	"io"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/xordb"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/rawdb"
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

// Send Transaction with Signature and txdata
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
		// Send transaction
		err := SendObject(conn, (*txs)[i])
		if err != nil {
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

// Send Interlinks Block
func SendInterlinks(conn net.Conn, interlinks []uint64, bc *core.BlockChain) error {
	log.Printf("INTERLINKS : %v\n", interlinks)
	// Send interlinkss length
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(interlinks)))
	if _, err := conn.Write(length); nil != err {
		log.Printf("failed to send interlinks length; err: %v", err)
		return err
	}

        for i := 0; i < len(interlinks); i++ {
                // Send block
                err := SendBlock(conn, bc.BlockAt(interlinks[i]))
                if err != nil {
                        return err
                }
        }

	return nil
}

// Send state map
func SendState(conn net.Conn, db xordb.Database) error {
	// Send state size
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(rawdb.CountStates(db)))
	if _, err := conn.Write(length); nil != err {
		log.Printf("failed to send state length; err: %v", err)
		return err
	}

	// Send state (address - txhash) and txs from DB
	iter := db.NewIterator()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if string(key[0]) == "s" { // prefix for state
			// Send address
			err := SendMessage(conn, key)
			if err != nil {
				return err
			}
			// Send tx hash
			err = SendMessage(conn, value)
			if err != nil {
				return err
			}
			// Send tx
			tx, _, _, _ := rawdb.ReadTransaction(db, common.BytesToHash(value))
			err = SendObject(conn, tx)
			if err != nil {
				return err
			}
		}
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

func RecvState(conn net.Conn, db xordb.Database) (error) {
	// Get State length
	statelen, err := RecvLength(conn)
	if err != nil {
		return err
	}

	// Get Address - txHash
	for i := uint32(0); i < statelen; i++ {
		// Get Address byte[]
		addrbuf, err := RecvObjectJson(conn)
		if err != nil {
			return err
		}

		// Get txHash byte[]
		txhashbuf, err := RecvObjectJson(conn)
		if err != nil {
			return err
		}

		// Insert to state map
		// [FIXME] Just Byte[] api would be OK
		rawdb.WriteState(db, common.BytesToAddress(addrbuf), common.BytesToHash(txhashbuf))

		// Get tx
		txbuf, err := RecvObjectJson(conn)
		if err != nil {
			return err
		}
		tx := types.UnmarshalJSON(txbuf)
	
		// Validate tx
		err = tx.ValidateTx()
		if err != nil {
			return err
		}

		// Write Transaction in DB
		rawdb.WriteTransaction(db, common.BytesToHash(txhashbuf), tx)
	}

	return nil
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
		// Get transaction
		txbuf, err := RecvObjectJson(conn)
		if err != nil {
			return nil, err
		}

		// Unmarshal Tx
		tx := types.UnmarshalJSON(txbuf)
		txs.Insert(tx)
	}

	// Make Block with header, txs
	block := types.NewBlock(&header, txs)
	block.Hash() // set block hash
	return block, nil
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number 
// of garage collection cycles completed.
func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

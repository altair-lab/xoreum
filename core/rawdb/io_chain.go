package rawdb

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/log"
	"github.com/altair-lab/xoreum/rlp"
	"github.com/altair-lab/xoreum/xordb"
)

// ReadHash retrieves the hash assigned to a block number.
func ReadHash(db xordb.Reader, number uint64) common.Hash {
	data, _ := db.Get(headerHashKey(number))
	fmt.Println("key:", common.Bytes2Hex(headerHashKey(number)))
	fmt.Println("hash value:", common.Bytes2Hex(data))
	if len(data) == 0 {
		fmt.Println("no", number)
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHash stores the hash assigned to a block number.
func WriteHash(db xordb.Writer, hash common.Hash, number uint64) {
	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		log.Crit("Failed to store number to hash mapping", "err", err)
		fmt.Println("hash not written")
	}
	fmt.Println("hash key:", common.Bytes2Hex(headerHashKey(number)))
	fmt.Println("hash value: ", hash.ToHex())
}

// DeleteHash removes the number to hash mapping.
func DeleteHash(db xordb.Writer, number uint64) {
	if err := db.Delete(headerHashKey(number)); err != nil {
		log.Crit("Failed to delete number to hash mapping", "err", err)
	}
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db xordb.Reader, hash common.Hash) *uint64 {
	data, _ := db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// ReadHeaderData retrieves a block header string
func ReadHeaderData(db xordb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	fmt.Println("header key:", common.Bytes2Hex(headerKey(number, hash)))
	data, _ := db.Get(headerKey(number, hash))

	fmt.Println("rlp data:", common.Bytes2Hex(data))
	return data
}

// HasHeader verifies the existence of a block header corresponding to the hash.
func HasHeader(db xordb.Reader, hash common.Hash, number uint64) bool {
	if has, err := db.Has(headerKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadHeader retrieves the block header corresponding to the hash.
func ReadHeader(db xordb.Reader, hash common.Hash, number uint64) *types.Header {
	data := ReadHeaderData(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		log.Error("Invalid block header", "hash", hash, "err", err)
		return nil
	}
	return header
}

// WriteHeader stores a block header into the database and also stores the hash-
// to-number mapping.
func WriteHeader(db xordb.Writer, header *types.Header) {
	fmt.Println("writing header")
	// Write the hash -> number mapping
	var (
		hash    = header.Hash()
		number  = header.Number
		encoded = encodeBlockNumber(number)
	)
	key := headerNumberKey(hash)
	if err := db.Put(key, encoded); err != nil {
		log.Crit("Failed to store hash to number mapping", "err", err)
	}
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		fmt.Println("header write fail")
		log.Crit("Failed to RLP encode header", "err", err)
	}

	key = headerKey(number, hash)
	fmt.Println("header key:", common.Bytes2Hex(key))
	fmt.Println("	======before encode======")
	header.PrintHeader()
	fmt.Println("	======after encode======")
	fmt.Println("	rlp data:", common.Bytes2Hex(data))

	if err := db.Put(key, data); err != nil {
		fmt.Println("header write fail")
		log.Crit("Failed to store header", "err", err)
	}
	fmt.Println("")
}

// DeleteHeader removes all block header data associated with a hash.
func DeleteHeader(db xordb.Writer, hash common.Hash, number uint64) {
	if err := db.Delete(headerNumberKey(hash)); err != nil {
		log.Crit("Failed to delete hash to number mapping", "err", err)
	}
}

// ReadBodyData retrieves the block body (transactions and uncles) in RLP encoding.
func ReadBodyData(db xordb.Reader, hash common.Hash, number uint64) []byte {
	data, _ := db.Get(blockBodyKey(number, hash))
	return data
}

// WriteBodyData stores block body into the database.
func WriteBodyData(db xordb.Writer, hash common.Hash, number uint64, data []byte) {
	fmt.Println("body key:", common.Bytes2Hex(blockBodyKey(number, hash)))
	fmt.Println("encoded data:", data)
	if err := db.Put(blockBodyKey(number, hash), data); err != nil {
		fmt.Println("store fail")
		log.Crit("Failed to store block body", "err", err)
	}
}

// // HasBody verifies the existence of a block body corresponding to the hash.
// func HasBody(db xordb.Reader, hash common.Hash, number uint64) bool {
// 	if has, err := db.Has(blockBodyKey(number, hash)); !has || err != nil {
// 		return false
// 	}
// 	return true
// }

// ReadBody retrieves the block body corresponding to the hash.
func ReadBody(db xordb.Reader, hash common.Hash, number uint64) *types.Body {
	data := ReadBodyData(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new(types.Body)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		log.Error("Invalid block body RLP", "hash", hash, "err", err)
		return nil
	}
	return body
}

// WriteBody stores a block body into the database.
func WriteBody(db xordb.Writer, hash common.Hash, number uint64, body *types.Body) {
	fmt.Println("writing body")

	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	dec := gob.NewDecoder(&data)

	fmt.Println(body.Transactions)
	err := enc.Encode(body)
	if err != nil {
		fmt.Println("error while encoding")
	}

	fmt.Println("encoded:", data.Bytes())

	var saved types.Body
	err = dec.Decode(&saved)
	if err != nil {
		log.Error("decode error:", err)
	}
	fmt.Println(saved.Transactions)
	saved.PrintBody()
	WriteBodyData(db, hash, number, data.Bytes())
}

// // DeleteBody removes all block body data associated with a hash.
// func DeleteBody(db xordb.Writer, hash common.Hash, number uint64) {
// 	if err := db.Delete(blockBodyKey(number, hash)); err != nil {
// 		log.Crit("Failed to delete block body", "err", err)
// 	}
// }

// ReadTdData retrieves a block's total difficulty corresponding to the hash.
func ReadTdData(db xordb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	data, _ := db.Get(headerTDKey(number, hash))
	return data
}

// ReadTd retrieves a block's total difficulty corresponding to the hash.
func ReadTd(db xordb.Reader, hash common.Hash, number uint64) *big.Int {
	data := ReadTdData(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	td := new(big.Int)
	if err := rlp.Decode(bytes.NewReader(data), td); err != nil {
		log.Error("Invalid block total difficulty RLP", "hash", hash, "err", err)
		return nil
	}
	return td
}

// WriteTd stores the total difficulty of a block into the database.
func WriteTd(db xordb.Writer, hash common.Hash, number uint64, td *big.Int) {
	data, err := rlp.EncodeToBytes(td)
	if err != nil {
		log.Crit("Failed to RLP encode block total difficulty", "err", err)
	}
	if err := db.Put(headerTDKey(number, hash), data); err != nil {
		log.Crit("Failed to store block total difficulty", "err", err)
	}
}

// DeleteTd removes all block total difficulty data associated with a hash.
func DeleteTd(db xordb.Writer, hash common.Hash, number uint64) {
	if err := db.Delete(headerTDKey(number, hash)); err != nil {
		log.Crit("Failed to delete block total difficulty", "err", err)
	}
}

// LoadBlockByBN retrieves an entire block corresponding to the number
func LoadBlockByBN(db xordb.Reader, number uint64) *types.Block {
	hash := ReadHash(db, number)

	header := ReadHeader(db, hash, number)

	if header == nil {
		fmt.Println("header empty")
		return nil
	}

	body := ReadBody(db, hash, number)

	txs := body.Transactions

	fmt.Println("body txs:", txs)

	b := types.NewBlock(header, txs)

	b.PrintTxs()

	return b
}

// LoadHeaderByBN retrieves an entire header corresponding to the number
func LoadHeaderByBN(db xordb.Reader, number uint64) *types.Header {
	hash := ReadHash(db, number)

	header := ReadHeader(db, hash, number)

	if header == nil {
		fmt.Println("header empty")
		return nil
	}

	return types.CopyHeader(header)
}

// LoadHeader retrieves an entire header corresponding to the hash & number
func LoadHeader(db xordb.Reader, hash common.Hash, number uint64) *types.Header {
	header := ReadHeader(db, hash, number)

	if header == nil {
		fmt.Println("header empty")
		return nil
	}

	return types.CopyHeader(header)
}

// LoadBlock retrieves an entire block corresponding to the hash & number
func LoadBlock(db xordb.Reader, hash common.Hash, number uint64) *types.Block {
	header := ReadHeader(db, hash, number)

	if header == nil {
		fmt.Println("header nil")
		return nil
	}
	//tx, blockHash, *blockNumber, uint64(txIndex)
	tx, _, _, _ := ReadTransaction(db, hash)
	if tx == nil {
		fmt.Println("tx nil")
		return nil
	}

	txs := []*types.Transaction{tx}
	return types.NewBlock(header, txs)
}

// StoreBlock serializes a block into the database, header and body separately.
func StoreBlock(db xordb.Writer, block *types.Block) {
	WriteHash(db, block.Hash(), block.Number())
	WriteHeader(db, block.Header())
	WriteBody(db, block.Hash(), block.Number(), block.Body())
}

// DeleteBlock removes all block data associated with a hash.
func DeleteBlock(db xordb.Writer, hash common.Hash, number uint64) {
	DeleteHash(db, number)
	DeleteHeader(db, hash, number)
	// DeleteBody(db, hash, number)
	DeleteTd(db, hash, number)
}

// ReadLastHeaderHash retrieves the hash of the current canonical head header.
func ReadLastHeaderHash(db xordb.Reader) common.Hash {
	data, _ := db.Get(lastHeaderKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteLastHeaderHash stores the hash of the current canonical head header.
func WriteLastHeaderHash(db xordb.Writer, hash common.Hash) {
	if err := db.Put(lastHeaderKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last header's hash", "err", err)
	}
}

// // ReadHeadBlockHash retrieves the hash of the current canonical head block.
// func ReadHeadBlockHash(db xordb.Reader) common.Hash {
// 	data, _ := db.Get(headBlockKey)
// 	if len(data) == 0 {
// 		return common.Hash{}
// 	}
// 	return common.BytesToHash(data)
// }

// // WriteHeadBlockHash stores the head block's hash.
// func WriteHeadBlockHash(db xordb.Writer, hash common.Hash) {
// 	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
// 		log.Crit("Failed to store last block's hash", "err", err)
// 	}
// }

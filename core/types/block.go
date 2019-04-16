package types

import (
	"fmt"
	"sync/atomic"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/crypto"
)

type Header struct {
	ParentHash common.Hash    `parentHash` // previous block's hash
	Coinbase   common.Address `miner`
	Root       common.Hash    `stateRoot`
	TxHash     common.Hash    `transactionsRoot`
	State      state.State    `state`
	Number     uint64         `number` // (Number == Height) A scalar value equal to the number of ancestor blocks. The genesis block has a number of zero
	Time       uint64         `timestamp`
	Nonce      uint64         `nonce`
	InterLink  []uint64       `interlink` // list of block's level
	Difficulty uint64         `difficulty`
}

/*
type Body struct {
	Transactions	[]*Transaction
	Uncles		[]*Header
}
*/

type Block struct {
	header       *Header
	transactions Transactions
	hash         atomic.Value
	level        uint64 // used in interlink
}

func (h *Header) Hash() common.Hash {
	return crypto.Keccak256Hash(common.ToBytes(*h))
}

// block's hash is same with header's hash
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

func (b *Block) PrintTx() {
	for i := 0; i < len(b.transactions); i++ {
		fmt.Println("====================")
		fmt.Println("Sender: ", b.transactions[i].Sender())
		fmt.Println("Recipient: ", b.transactions[i].Recipient())
		fmt.Println("Value: ", b.transactions[i].Value())
	}
}

func (b *Block) InsertTx(tx *Transaction) {
	b.transactions = append(b.transactions, tx)
}

func NewBlock(header *Header, txs []*Transaction) *Block {
	return &Block{
		header:       header,
		transactions: txs,
	}
}

func NewHeader(parentHash common.Hash, miner common.Address, stateRoot common.Hash, txHash common.Hash, state state.State, difficulty uint64, time uint64, nonce uint64) *Header {
	return &Header{
		ParentHash: parentHash,
		Coinbase:   miner,
		Root:       stateRoot,
		TxHash:     txHash,
		State:      state,
		Difficulty: difficulty,
		Time:       time,
		Nonce:      nonce,
	}
}

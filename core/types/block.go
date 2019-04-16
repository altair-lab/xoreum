package types

import (
	"sync/atomic"
	"fmt"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/crypto"
)

type Header struct {
	ParentHash common.Hash    `parentHash`
	Coinbase   common.Address `miner`
	Root       common.Hash    `stateRoot`
	TxHash     common.Hash    `transactionsRoot`
	State      state.State    `state`
	Difficulty uint64         `difficulty`
	Time       uint64         `timestamp`
	Nonce      uint64         `nonce`
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

func (b* Block) PrintTx() {
	for i := 0; i < len(b.transactions); i++ {
		fmt.Printf("Sender: %x\n", b.transactions[i].Sender())
		fmt.Printf("Recipient: %x\n", b.transactions[i].Recipient())
		fmt.Printf("Value: %d\n", b.transactions[i].Value())
		fmt.Printf("====================\n")
	}
}

func (b* Block) InsertTx(tx *Transaction) {
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

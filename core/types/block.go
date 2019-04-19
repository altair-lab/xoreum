package types

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/crypto"
)

const (
	InterlinkLength = 10
)

type Header struct {
	ParentHash common.Hash             `parentHash` // previous block's hash
	Coinbase   common.Address          `miner`
	Root       common.Hash             `stateRoot`
	TxHash     common.Hash             `transactionsRoot`
	State      state.State             `state`
	Number     uint64                  `number` // (Number == Height) A scalar value equal to the number of ancestor blocks. The genesis block has a number of zero
	Time       uint64                  `timestamp`
	Nonce      uint64                  `nonce`
	InterLink  [InterlinkLength]uint64 `interlink` // list of block's level
	Difficulty uint64                  `difficulty`
}

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

func (b *Block) GetLevel() uint64 {
	var level uint64 = 0
	//dif := core.Difficulty
	//dif := big.NewInt(10000)
	dif := common.Difficulty
	blockHash := b.Hash().ToBigInt()

	for {
		// if blockhash < dif
		if blockHash.Cmp(dif) == -1 {
			dif = new(big.Int).Div(dif, big.NewInt(2)) // dif /= 2
			level++
		} else {
			break
		}
	}

	// level starts from 0
	// set block's level
	b.level = level - 1

	return b.level
}

// return interlink that contains this block too
// to compare with next block's interlink
// should be current_block.GetUpdatedInterlink() == next_block.header.Interlink
// Also, this function can be used when you fill newly mined block's interlink
// new_mined_block.header.Interlink = current_block.GetUpdatedInterlink()
func (b *Block) GetUpdatedInterlink() [InterlinkLength]uint64 {
	// copy interlink
	updatedInterlink := b.header.InterLink

	// get updated interlink
	lv := b.GetLevel()
	if lv > 10 {
		lv = 10
	}
	for i := uint64(0); i < lv; i++ {
		updatedInterlink[i] = b.header.Number
	}

	return updatedInterlink
}

func (b *Block) GetHeader() *Header {
	return b.header
}

func (b *Block) PrintBlock() {
	fmt.Println("====================")
	fmt.Println("block number:", b.header.Number)
	fmt.Println("block parent hash:", b.header.ParentHash.ToHex())
	fmt.Println("       block hash:", b.Hash().ToHex())
	fmt.Println("block level:", b.GetLevel())
	fmt.Println("block nonce:", b.header.Nonce)
	fmt.Println("block interlink:", b.header.InterLink)
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

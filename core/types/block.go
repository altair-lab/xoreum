package types

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/crypto"
)

const (
	InterlinkLength = uint64(10)
)

type Header struct {
	ParentHash common.Hash             `json:"parentHash"` // previous block's hash
	Coinbase   common.Address          `json:"miner"`
	Root       common.Hash             `json:"stateRoot"`
	TxHash     common.Hash             `json:"transactionsRoot"`
	Number     uint64                  `json:"number"` // (Number == Height) A scalar value equal to the number of ancestor blocks. The genesis block has a number of zero
	Time       uint64                  `json:"timestamp"`
	Nonce      uint64                  `json:"nonce"`
	InterLink  [InterlinkLength]uint64 `json:"interlink"`  // list of block's level
	Difficulty uint64                  `json:"difficulty"` // no difficulty change, so set global Difficulty
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

func (b *Block) PrintTxs() {
	for i := 0; i < len(b.transactions); i++ {
		fmt.Println("====================")
		fmt.Println("tx ", i)
		b.transactions[i].PrintTx()
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
	if lv > InterlinkLength {
		lv = InterlinkLength
	}
	for i := uint64(0); i < lv; i++ {
		updatedInterlink[i] = b.header.Number
	}

	return updatedInterlink
}

func (b *Block) GetUniqueInterlink() []uint64 {
	// include current block
	list := unique(b.header.InterLink)
	list = append(list, b.GetHeader().Number)
	return list
}

func unique(intSlice [InterlinkLength]uint64) []uint64 {
	keys := make(map[uint64]bool)
	list := []uint64{}
	for i := len(intSlice) - 1; i >= 0; i-- {
		if _, value := keys[intSlice[i]]; !value {
			keys[intSlice[i]] = true
			list = append(list, intSlice[i])
		}
	}
	return list
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
	b.PrintTxs()
}

func NewBlock(header *Header, txs []*Transaction) *Block {
	return &Block{
		header:       header,
		transactions: txs,
	}
}

func NewHeader(parentHash common.Hash, miner common.Address, stateRoot common.Hash, txHash common.Hash, difficulty uint64, number uint64, time uint64, nonce uint64) *Header {
	return &Header{
		ParentHash: parentHash,
		Coinbase:   miner,
		Root:       stateRoot,
		TxHash:     txHash,
		Difficulty: difficulty,
		Number:     number,
		Time:       time,
		Nonce:      nonce,
	}
}

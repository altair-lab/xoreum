package types

import (
	"fmt"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/rlp"
	"golang.org/x/crypto/sha3"
)

var (
	EmptyRootHash  = DeriveSha(Transactions{})
	EmptyUncleHash = rlpHash([]*Header(nil))
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

type Body struct {
	Transactions []*Transaction
}

type Block struct {
	header       *Header
	transactions Transactions
	hash         atomic.Value
	level        uint64 // used in interlink

	size atomic.Value
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

// GetUpdatedInterlink returns interlink that contains this block too
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

func (b *Block) GetTxs() *Transactions {
	return &b.transactions
}

func (b *Block) PrintBlock() {
	fmt.Println("====================")
	fmt.Println("block number:", b.header.Number)
	fmt.Println("block parent hash:", b.header.ParentHash.ToHex())
	fmt.Println("block hash:", b.Hash().ToHex())
	fmt.Println("block level:", b.GetLevel())
	fmt.Println("block nonce:", b.header.Nonce)
	fmt.Println("block time:", b.header.Time)
	fmt.Println("block interlink:", b.header.InterLink)
	b.PrintTxs()
}

func (h *Header) PrintHeader() {
	fmt.Println("	[HEADER]")
	fmt.Println("	hash:", h.Hash().ToHex())
	fmt.Println("	parent hash:", h.ParentHash.ToHex())
	fmt.Println("	number:", h.Number)
	fmt.Println("	interlink:", h.InterLink)
}

func NewBlock(header *Header, txs []*Transaction) *Block {
	// b := &Block{header: CopyHeader(header)}
	// if len(txs) == 0 {
	// 	b.header.TxHash = EmptyRootHash
	// } else {
	// 	b.header.TxHash = DeriveSha(Transactions(txs))
	// 	b.transactions = make(Transactions, len(txs))
	// 	copy(b.transactions, txs)

	// }
	// return &Block{
	// 	header:       CopyHeader(header),
	// 	transactions: txs,
	// }
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
func CopyHeader(header *Header) *Header {
	return &Header{
		ParentHash: header.ParentHash,
		Coinbase:   header.Coinbase,
		Root:       header.Root,
		TxHash:     header.TxHash,
		Difficulty: header.Difficulty,
		Number:     header.Number,
		Time:       header.Time,
		Nonce:      header.Nonce,
	}
}

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header *Header
	Txs    []*Transaction
}

// DecodeRLP decodes the Ethereum
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.transactions = eb.Header, eb.Txs
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		Header: b.header,
		Txs:    b.transactions,
	})
}

func (b *Block) Number() uint64 { return (b.header.Number) }

func (b *Block) Header() *Header { return (b.header) }

func (b *Block) Body() *Body { return &Body{b.transactions} }

func (b *Block) Transactions() Transactions { return b.transactions }

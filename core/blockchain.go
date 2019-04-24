package core

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/params"
)

var (
	// incorrect block's number (current_block_number + 1 != insert_block's_number)
	ErrWrongBlockNumber = errors.New("incorrect block number")

	ErrWrongParentHash = errors.New("block's parent hash does not match with current block")

	ErrTooHighHash = errors.New("block's hash is higher than difficulty")

	ErrWrongInterlink = errors.New("wrong interlink")
)

type BlockChain struct {
	//ChainID *big.Int // chainId identifies the current chain and is used for replay protection

	genesisBlock *types.Block
	currentBlock atomic.Value
	//processor	Processor
	//validator Validator

	blocks []types.Block // temporary block list. blocks will be saved in db
}

func NewBlockChain() *BlockChain {

	bc := &BlockChain{
		//ChainID:      big.NewInt(0),
		genesisBlock: params.GetGenesisBlock(),
	}
	bc.currentBlock.Store(bc.genesisBlock)
	bc.blocks = append(bc.blocks, *bc.genesisBlock)

	return bc
}

// check block's validity, if ok, then insert block into chain
func (bc *BlockChain) Insert(block *types.Block) error {

	// validate block before insert
	err := bc.validateBlock(block)

	if err != nil {
		// didn't pass validation
		return err
	} else {
		// pass all validation
		// insert that block into blockchain
		bc.insert(block)
		return nil
	}
}

// check that this block is valid to be inserted
func (bc *BlockChain) validateBlock(block *types.Block) error {

	// 1. check block number
	if bc.CurrentBlock().GetHeader().Number+1 != block.GetHeader().Number {
		return ErrWrongBlockNumber
	}

	// 2. check parent hash
	if bc.CurrentBlock().Hash() != block.GetHeader().ParentHash {
		return ErrWrongParentHash
	}

	// 3. check that block hash < difficulty
	if block.GetHeader().Hash().ToBigInt().Cmp(common.Difficulty) != -1 {
		return ErrTooHighHash
	}

	// 4. check block's interlink
	if bc.CurrentBlock().GetUpdatedInterlink() != block.GetHeader().InterLink {
		return ErrWrongInterlink
	}

	// 5. check trie

	// 6. check txs

	// pass all validation. return no err
	return nil
}

// actually insert block
func (bc *BlockChain) insert(block *types.Block) {
	bc.blocks = append(bc.blocks, *block)
	bc.currentBlock.Store(block)
}

func (bc *BlockChain) CurrentBlock() *types.Block {
	return bc.currentBlock.Load().(*types.Block)
}

func (bc *BlockChain) PrintBlockChain() {
	fmt.Println("=== Print Blocks ===")
	for i := 0; i < len(bc.blocks); i++ {
		bc.blocks[i].PrintBlock()
	}
	fmt.Println("====================")
	fmt.Println("=== End of Chain ===")
}

// make blockchain for test. insert simple blocks
func MakeTestBlockChain(chainLength uint64) *BlockChain {

	bc := NewBlockChain()

	var empty_txs []*types.Transaction
	empty_txs = []*types.Transaction{}

	// insert blocks into blockchain
	for i := uint64(1); i <= chainLength; i++ {
		b := types.NewBlock(&types.Header{}, empty_txs)
		b.GetHeader().ParentHash = bc.CurrentBlock().Hash()
		b.GetHeader().Number = i
		b.GetHeader().Nonce = 0
		b.GetHeader().InterLink = bc.CurrentBlock().GetUpdatedInterlink()

		// simple PoW
		for {
			err := bc.Insert(b)

			if err != nil {
				b.GetHeader().Nonce++
			} else {
				break
			}
		}
	}

	return bc
}

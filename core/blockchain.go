package core

import (
	"sync/atomic"

	"github.com/altair-lab/xoreum/common/math"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/params"
)

var (
	Difficulty = math.BigPow(2, 256-10) // mining difficulty: 100
)

type BlockChain struct {
	//ChainID *big.Int // chainId identifies the current chain and is used for replay protection

	genesisBlock *types.Block
	currentBlock atomic.Value
	//processor	Processor
	//validator Validator

	blocks []types.Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		//ChainID:      big.NewInt(0),
		genesisBlock: params.GetGenesisBlock(),
	}
}

func (bc *BlockChain) insert(block *types.Block) {

	bc.blocks = append(bc.blocks, *block)
	bc.currentBlock.Store(block)
}

func (bc *BlockChain) CurrentBlock() *types.Block {
	return bc.currentBlock.Load().(*types.Block)
}

package core

import (
	"sync/atomic"

	"github.com/altair-lab/xoreum/core/types"
)

type BlockChain struct {
	genesisBlock *types.Block
	currentBlock atomic.Value
	//processor	Processor
	validator Validator

	blocks []types.Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{}
}

func (bc *BlockChain) insert(block *types.Block) {

	bc.blocks = append(bc.blocks, *block)
	bc.currentBlock.Store(block)
}

func (bc *BlockChain) CurrentBlock() *types.Block {
	return bc.currentBlock.Load().(*types.Block)
}

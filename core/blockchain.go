package core

import (
	"github.com/altair-lab/xoreum/core/types"
)

type BlockChain struct {
	genesisBlock	*types.Block
	currentBlock	*types.Block
	//processor	Processor
	//validator	Validator

	blocks		[]types.Block
}


func (bc *BlockChain) insert(block *types.Block){
	
	bc.blocks = append(bc.blocks, *block);
	bc.currentBlock = block;
}




package core

import (
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/consensus"
)

type BlockValidator struct {
	bc *BlockChain		// Canonical blockchain
	engine consensus.Engine	// Consensus engine used for validating
}

func (v *BlockValidator) ValidateBody(block *types.Block) error {

	if v.bc.CurrentBlock().Hash() == block.Hash() {
		return ErrKnownBlock
	}

	return nil
}

func (v *BlockValidator) ValidateState(block *types.Block) error {
	return nil
}

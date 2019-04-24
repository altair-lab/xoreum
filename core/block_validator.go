package core

import (
	"github.com/altair-lab/xoreum/consensus"
	"github.com/altair-lab/xoreum/core/types"
)

type BlockValidator struct {
	bc     *BlockChain      // Canonical blockchain
	engine consensus.Engine // Consensus engine used for validating
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

package core

import (
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
)

type StateProcessor struct {
	bc *BlockChain
}

func (p *StateProcessor) Process(block *types.Block, state *state.State) {

}

func ApplyTransaction(tx *types.Transaction) {

}

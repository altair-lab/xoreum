package core

import(
	"github.com/altair-lab/xoreum/core/types"
)

type BlockValidator struct{

	bc	*BlockChain
}

func (v *BlockValidator) ValidateBody(block *types.Block) error{
	return nil
}

func (v *BlockValidator) ValidateState(block *types.Block) error{
	return nil
}



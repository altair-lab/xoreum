package core

import (
	"fmt"
	"github.com/altair-lab/xoreum/core/types"
)

type BlockChain struct {
	genesisBlock	*types.Block
	currentBlock	*types.Block
	//processor	Processor
	//validator	Validator
}

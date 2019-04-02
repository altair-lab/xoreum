package core

import (
	"fmt"
	"github.com/altair-lab/xoreum/core/types"
)

type BlockChain struct {
	genesisBlock	*Block
	currentBlock	*Block
	//processor	Processor
	//validator	Validator
}

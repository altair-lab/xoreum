package main

import (
	"fmt"
	"types"
)

type BlockChain struct {
	genesisBlock	*Block
	currentBlock	*Block
	//processor	Processor
	//validator	Validator
}

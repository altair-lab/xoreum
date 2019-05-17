package test

import (
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/xordb/memorydb"
)

func ExampleFunc6() {

	// copy bitcoin's genesis block
	db := memorydb.New()
	bc := core.NewBlockChainForBitcoin(db)
	bc.PrintBlockChain()
	bc.GetState().Print()

	// output: 1
}

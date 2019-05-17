package test

import (
	"fmt"

	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/xordb/memorydb"
)

func ExampleFunc6() {

	// copy bitcoin's genesis block
	db := memorydb.New()
	bc, gpriv := core.NewBlockChainForBitcoin(db)
	bc.PrintBlockChain()
	bc.GetState().Print()

	fmt.Println("\ngenesis account's private key:", gpriv)

	// output: 1
}

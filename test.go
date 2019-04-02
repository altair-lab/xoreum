package main

import (
	"fmt"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/common"
)

func main() {

	fmt.Println("---test common/types.go---")
	hl := common.HashLength
	al := common.AddressLength
	//hash1 := common.Hash
	//address1 := common.Address

	fmt.Println("---test core---")
	blockchain1 := core.BlockChain{}

	fmt.Println("---test core/state---")
	state1 := state.State
	account1 := state.Account{}

	fmt.Println("---test core/types---")
	header1 := types.Header{}
	block1 := types.Block{}
	transaction1 := types.Transaction{}

	fmt.Println("---all test passed---")

}

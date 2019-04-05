package main

import (
	"fmt"
	"math/rand"

	"github.com/altair-lab/xoreum/common"
	//"github.com/altair-lab/xoreum/common/hexutil"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
)

func makeTestBlockChain(){
	
	return	

}

func main() {

	fmt.Println("---test common/types.go---")
	hl := common.HashLength
	al := common.AddressLength
	hash1 := common.Hash{}
	hash2 := make([]byte, 32)
	hash3 := crypto.Keccak256Hash()
	rand.Read(hash2)
	address1 := common.Address{}
	fmt.Println(hl)
	fmt.Println(al)
	fmt.Println(hash1)
	fmt.Println(address1)

	fmt.Println("---test core---")
	blockchain1 := core.BlockChain{}
	fmt.Println(blockchain1)

	fmt.Println("---test core/state---")
	state1 := state.State{}
	account1 := state.Account{}
	fmt.Println(state1)
	fmt.Println(account1)

	fmt.Println("---test core/types---")
	header1 := types.Header{}
	block1 := types.Block{}
	transaction1 := types.Transaction{}
	fmt.Println(header1)
	fmt.Println(block1)
	fmt.Println(transaction1)

	fmt.Println("---test core/crypto---")
	keccak256_1 := crypto.Keccak256()
	fmt.Println(keccak256_1)

	fmt.Println("---all test passed---")

	fmt.Println("hash2: ", hash3.ToHex())


	fmt.Println("\n\n\n\n\n\n")



}

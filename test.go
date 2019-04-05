package main

import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
	//"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/miner"
	//"github.com/altair-lab/xoreum/crypto"
)

func main() {
	// send "10" to "account 1"
	
	// Create default accounts
	fmt.Println("===== Create default accounts . . . =====")
	acc0 := state.NewAccount(common.Address{0}, uint64(0), uint64(100)) // acc0 [Address:0, Nonce:0, Balance:100]
	acc1 := state.NewAccount(common.Address{1}, uint64(0), uint64(100)) // acc1 [Address:1, Nonce:0, Balance:100]
	acc2 := state.NewAccount(common.Address{2}, uint64(0), uint64(0)) // acc1 [Address:1, Nonce:0, Balance:0]
	fmt.Println("Account0: ", acc0)
	fmt.Println("Account1: ", acc1)
	fmt.Println("Account2: ", acc2)

	// [TODO]Set state
	state := state.State{}

	// Create transaction (send "10" from "account0" to "account1")
	fmt.Println("===== Create transaction that send [10] from [account0] to [account1] =====")
	tx := types.NewTransaction(acc0.Address, acc1.Address, uint64(10))
	fmt.Println("Transaction: ", tx)

	// Mining
	fmt.Println("===== Start Mining (Miner : account2) =====")
	miner := miner.Miner{acc2.Address}
	block := miner.Mine(*tx, state)
	fmt.Println(block)



	/*
	fmt.Println("---test common/types.go---")
	hl := common.HashLength
	al := common.AddressLength
	hash1 := common.Hash{}
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
*/
}

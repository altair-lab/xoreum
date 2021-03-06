package miner 

/*
import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
	//"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/miner"
	//"github.com/altair-lab/xoreum/crypto"
)

func ExampleMining() {
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
	block := miner.Mine(*tx, state, 240) // Difficulty is 0~255 for now
	fmt.Println("Created Block: " , block)

	// output:
	// "a"
}
*/

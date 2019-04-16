package test

/*
import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/miner"
)



func ExampleTxpool() {

	// Make account
	acc0 := state.NewAccount(common.Address{0}, uint64(0), uint64(7000)) // acc0 [Address:0, Nonce:0, Balance:7000]
	acc1 := state.NewAccount(common.Address{1}, uint64(0), uint64(2000)) // acc1 [Address:1, Nonce:0, Balance:2000]
	acc2 := state.NewAccount(common.Address{2}, uint64(0), uint64(100))  // acc2 [Address:2, Nonce:0, Balance:100]
	acc3 := state.NewAccount(common.Address{3}, uint64(0), uint64(100))  // acc3 [Address:3, Nonce:0, Balance:100]

	// [TODO] Set state
	state := state.State{}

	// Create tranaction
	tx0 := types.NewTransaction(acc0.Address, acc1.Address, uint64(2000)) // send [2000] from [account0] to [account1]
	tx1 := types.NewTransaction(acc1.Address, acc0.Address, uint64(500)) // send [500] from [account1] to [account0]
	tx_overwrite_invalid := types.NewTransaction(acc0.Address, acc1.Address, uint64(100)) // send [100] from [account0] to [account1]
	tx_insufficient_invalid := types.NewTransaction(acc3.Address, acc2.Address, uint64(200)) // send [5000] from [account3] to [account2]
	// tx_nonce_invalid := types.NewTransaction(acc0.Address, acc1.Address, uint64(5000)) // send [5000] from [account0] to [account1]
	// tx_sender_invalid := types.NewTransaction(acc0.Address, acc1.Address, uint64(5000)) // send [5000] from [account0] to [account1]

	// Create txpool
	txpool := core.NewTxPool()

	// Add txs to txpool
	success, err := txpool.Add(tx0)
	if !success {
		fmt.Println(err)
	}

	success, err = txpool.Add(tx1)
	if !success {
		fmt.Println(err)
	}

	success, err = txpool.Add(tx_overwrite_invalid)
	if !success {
		fmt.Println(err)
	}

	success, err = txpool.Add(tx_insufficient_invalid)
	if !success {
		fmt.Println(err)
	}


	// Mining from txpool
	miner := miner.Miner{acc0.Address}
	block := miner.Mine(txpool, state, 240)
	if block != nil {
		block.PrintTx()
	} else {
		fmt.Println("Mining Fail")
	}

	// output:
	// true
}
*/

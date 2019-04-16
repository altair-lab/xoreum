package test

import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/miner"
)

///////////////////////
///*   CHECKLIST   *///
///* 1. FIFO       *///
///* 2. ValidateTx *///
///////////////////////

func ExampleTxpool() {
	
	// Make keys
	privatekey0, _ := crypto.GenerateKey()
	publickey0 := privatekey0.PublicKey
	address0 := crypto.Keccak256Address(common.ToBytes(publickey0))
	fmt.Printf("publickey 0 : %x \n", address0)

	privatekey1, _ := crypto.GenerateKey()
	publickey1 := privatekey1.PublicKey
	address1 := crypto.Keccak256Address(common.ToBytes(publickey1))
	fmt.Printf("publickey 1 : %x \n", address1)

	privatekey2, _ := crypto.GenerateKey()
	publickey2 := privatekey2.PublicKey
	address2 := crypto.Keccak256Address(common.ToBytes(publickey2))
	fmt.Printf("publickey 2 : %x \n", address2)
	
	privatekey3, _ := crypto.GenerateKey()
	publickey3 := privatekey3.PublicKey
	address3 := crypto.Keccak256Address(common.ToBytes(publickey3))
	fmt.Printf("publickey 3 : %x \n", address3)
	
	// Make account
	acc0 := state.NewAccount(address0, uint64(0), uint64(7000)) // acc0 [Address:0, Nonce:0, Balance:7000]
	state.NewAccount(address1, uint64(0), uint64(2000)) // acc1 [Address:1, Nonce:0, Balance:2000]
	state.NewAccount(address2, uint64(0), uint64(100))  // acc2 [Address:2, Nonce:0, Balance:100]
	state.NewAccount(address3, uint64(0), uint64(100))  // acc3 [Address:3, Nonce:0, Balance:100]

	// [TODO] Set state
	state := state.State{}

	// Create tranaction
	tx0 := types.NewTransaction(publickey0, publickey1, uint64(2000)) // send [2000] from [account0] to [account1]
	tx1 := types.NewTransaction(publickey1, publickey0, uint64(500)) // send [500] from [account1] to [account0]
	tx_overwrite_invalid := types.NewTransaction(publickey0, publickey1, uint64(100)) // send [100] from [account0] to [account1]
	tx_insufficient_invalid := types.NewTransaction(publickey3, publickey2, uint64(200)) // send [5000] from [account3] to [account2]
	// tx_nonce_invalid := types.NewTransaction(publickey0, publickey1, uint64(5000)) // send [5000] from [account0] to [account1]
	// tx_sender_invalid := types.NewTransaction(publickey0, publickey1, uint64(5000)) // send [5000] from [account0] to [account1]

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

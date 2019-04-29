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
	fmt.Println("========== Create Accounts ==========")
	privatekey0, _ := crypto.GenerateKey()
	publickey0 := privatekey0.PublicKey
	address0 := crypto.Keccak256Address(common.ToBytes(publickey0))

	privatekey1, _ := crypto.GenerateKey()
	publickey1 := privatekey1.PublicKey
	address1 := crypto.Keccak256Address(common.ToBytes(publickey1))

	privatekey2, _ := crypto.GenerateKey()
	publickey2 := privatekey2.PublicKey
	address2 := crypto.Keccak256Address(common.ToBytes(publickey2))
	
	privatekey3, _ := crypto.GenerateKey()
	publickey3 := privatekey3.PublicKey
	address3 := crypto.Keccak256Address(common.ToBytes(publickey3))
	
	// Make account
	acc0 := state.NewAccount(address0, uint64(0), uint64(7000)) // acc0 [Address:0, Nonce:0, Balance:7000]
	acc1 := state.NewAccount(address1, uint64(0), uint64(2000)) // acc1 [Address:1, Nonce:0, Balance:2000]
	acc2 := state.NewAccount(address2, uint64(0), uint64(100))  // acc2 [Address:2, Nonce:0, Balance:100]
	acc3 := state.NewAccount(address3, uint64(0), uint64(100))  // acc3 [Address:3, Nonce:0, Balance:100]

	// [TODO] Set state
	state := state.NewState()
	state.Add(acc0)
	state.Add(acc1)
	state.Add(acc2)
	state.Add(acc3)
	state.Print()
	
	fmt.Printf("\n")

	// Create tranaction
	fmt.Println("========== Create Transactions ==========")
	tx0 := types.NewTransaction(0, publickey0, publickey1, uint64(2000)) // send [2000] from [account0] to [account1]
	tx1 := types.NewTransaction(0, publickey1, publickey0, uint64(500)) // send [500] from [account1] to [account0]
	tx_overwrite_invalid := types.NewTransaction(0, publickey0, publickey1, uint64(100)) // send [100] from [account0] to [account1]
	tx2 := types.NewTransaction(1, publickey0, publickey1, uint64(300)) // send [100] from [account0] to [account1]
	tx_insufficient_invalid := types.NewTransaction(0, publickey3, publickey2, uint64(200)) // send [5000] from [account3] to [account2]
	// tx_nonce_invalid := types.NewTransaction(0, publickey0, publickey1, uint64(5000)) // send [5000] from [account0] to [account1]
	// tx_sender_invalid := types.NewTransaction(0, publickey0, publickey1, uint64(5000)) // send [5000] from [account0] to [account1]

	// Sign to transaction
	tx0_signed, _ := types.SignTx(tx0, privatekey0) // sign by sender
	tx1_signed, _ := types.SignTx(tx1, privatekey1) // sign by sender
	tx_overwrite_signed, _ := types.SignTx(tx_overwrite_invalid, privatekey0) // sign by sender
	tx2_signed, _ := types.SignTx(tx2, privatekey0) // sign by sender
	tx_insufficient_signed,_ := types.SignTx(tx_insufficient_invalid, privatekey3) // sign by sender
	
	tx0_signed, _ = types.SignTx(tx0_signed, privatekey1) // sign by receiver
	tx1_signed, _ = types.SignTx(tx1_signed, privatekey0) // sign by receiver
	tx_overwrite_signed, _ = types.SignTx(tx_overwrite_signed, privatekey1) // sign by receiver
	tx2_signed, _ = types.SignTx(tx2_signed, privatekey1) // sign by receiver
	tx_insufficient_signed,_ = types.SignTx(tx_insufficient_signed, privatekey2) // sign by receiver

	// Create Chain, txpool
	bc := core.NewBlockChain()
	txpool := core.NewTxPool(state, bc)

	// Add txs to txpool
	success, err := txpool.Add(tx0_signed)
	if !success {
		fmt.Println(err)
	}

	success, err = txpool.Add(tx1_signed)
	if !success {
		fmt.Println(err)
	}

	success, err = txpool.Add(tx_overwrite_signed)
	if !success {
		fmt.Println(err)
	}

	success, err = txpool.Add(tx2_signed)
	if !success {
		fmt.Println(err)
	}

	success, err = txpool.Add(tx_insufficient_signed)
	if !success {
		fmt.Println(err)
	}

	fmt.Printf("\n")
	
	// Mining from txpool
	fmt.Println("============ Mining block  ============")
	miner := miner.Miner{acc0.Address}
	block := miner.Mine(txpool, 240)
	if block != nil {
		block.PrintTx()
	} else {
		fmt.Println("Mining Fail")
	}

	// Add to Blockchain
	err = bc.Insert(block)
	if err != nil {
		fmt.Println(err)
	}
	bc.PrintBlockChain()

	// output:
	// true
}
